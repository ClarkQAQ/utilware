// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cache

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"sync"

	"utilware/dep/x/tools/go/analysis"
	"utilware/dep/x/tools/go/packages"
	"utilware/dep/x/tools/internal/lsp/source"
	"utilware/dep/x/tools/internal/lsp/telemetry"
	"utilware/dep/x/tools/internal/span"
	"utilware/dep/x/tools/internal/telemetry/log"
	errors "utilware/dep/x/xerrors"
)

type snapshot struct {
	id   uint64
	view *view

	// mu guards all of the maps in the snapshot.
	mu sync.Mutex

	// ids maps file URIs to package IDs.
	// It may be invalidated on calls to go/packages.
	ids map[span.URI][]packageID

	// metadata maps file IDs to their associated metadata.
	// It may invalidated on calls to go/packages.
	metadata map[packageID]*metadata

	// importedBy maps package IDs to the list of packages that import them.
	importedBy map[packageID][]packageID

	// files maps file URIs to their corresponding FileHandles.
	// It may invalidated when a file's content changes.
	files map[span.URI]source.FileHandle

	// packages maps a packageKey to a set of packageHandles to which that file belongs.
	// It may be invalidated when a file's content changes.
	packages map[packageKey]*packageHandle

	// actions maps an actionkey to its actionHandle.
	actions map[actionKey]*actionHandle

	// workspacePackages contains the workspace's packages, which are loaded
	// when the view is created.
	workspacePackages map[packageID]packagePath

	// unloadableFiles keeps track of files that we've failed to load.
	unloadableFiles map[span.URI]struct{}
}

type packageKey struct {
	mode source.ParseMode
	id   packageID
}

type actionKey struct {
	pkg      packageKey
	analyzer *analysis.Analyzer
}

func (s *snapshot) ID() uint64 {
	return s.id
}

func (s *snapshot) View() source.View {
	return s.view
}

// Config returns the configuration used for the snapshot's interaction with the
// go/packages API.
func (s *snapshot) Config(ctx context.Context) *packages.Config {
	env, buildFlags := s.view.env()
	cfg := &packages.Config{
		Env:        env,
		Dir:        s.view.folder.Filename(),
		Context:    ctx,
		BuildFlags: buildFlags,
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedTypesSizes,
		Fset:    s.view.session.cache.fset,
		Overlay: s.buildOverlay(),
		ParseFile: func(*token.FileSet, string, []byte) (*ast.File, error) {
			panic("go/packages must not be used to parse files")
		},
		Logf: func(format string, args ...interface{}) {
			if s.view.options.VerboseOutput {
				log.Print(ctx, fmt.Sprintf(format, args...))
			}
		},
		Tests: true,
	}
	return cfg
}

func (s *snapshot) buildOverlay() map[string][]byte {
	s.mu.Lock()
	defer s.mu.Unlock()

	overlays := make(map[string][]byte)
	for uri, fh := range s.files {
		overlay, ok := fh.(*overlay)
		if !ok {
			continue
		}
		if overlay.saved {
			continue
		}
		// TODO(rstambler): Make sure not to send overlays outside of the current view.
		overlays[uri.Filename()] = overlay.text
	}
	return overlays
}

func (s *snapshot) PackageHandles(ctx context.Context, fh source.FileHandle) ([]source.PackageHandle, error) {
	if fh.Identity().Kind != source.Go {
		panic(fmt.Sprintf("called PackageHandles on a non-Go FileHandle"))
	}

	ctx = telemetry.File.With(ctx, fh.Identity().URI)

	// Check if we should reload metadata for the file. We don't invalidate IDs
	// (though we should), so the IDs will be a better source of truth than the
	// metadata. If there are no IDs for the file, then we should also reload.
	ids := s.getIDsForURI(fh.Identity().URI)
	reload := len(ids) == 0
	for _, id := range ids {
		// Reload package metadata if any of the metadata has missing
		// dependencies, in case something has changed since the last time we
		// reloaded it.
		if m := s.getMetadata(id); m == nil || len(m.missingDeps) > 0 {
			reload = true
			break
		}
	}
	if reload {
		if err := s.load(ctx, fileURI(fh.Identity().URI)); err != nil {
			return nil, err
		}
	}
	// Get the list of IDs from the snapshot again, in case it has changed.
	var phs []source.PackageHandle
	for _, id := range s.getIDsForURI(fh.Identity().URI) {
		ph, err := s.packageHandle(ctx, id, source.ParseFull)
		if err != nil {
			return nil, err
		}
		phs = append(phs, ph)
	}
	return phs, nil
}

// packageHandle returns a PackageHandle for the given ID. It assumes that
// the metadata for the given ID has already been loaded, but if the
// PackageHandle has not been constructed, it will rebuild it.
func (s *snapshot) packageHandle(ctx context.Context, id packageID, mode source.ParseMode) (*packageHandle, error) {
	ph := s.getPackage(id, mode)
	if ph != nil {
		return ph, nil
	}
	// Don't reload metadata in this function.
	// Callers of this function must reload metadata themselves.
	m := s.getMetadata(id)
	if m == nil {
		return nil, errors.Errorf("%s has no metadata", id)
	}
	return s.buildPackageHandle(ctx, m.id, mode)
}

func (s *snapshot) GetReverseDependencies(ctx context.Context, id string) ([]source.PackageHandle, error) {
	if err := s.awaitLoaded(ctx); err != nil {
		return nil, err
	}
	ids := make(map[packageID]struct{})
	s.transitiveReverseDependencies(packageID(id), ids)

	// Make sure to delete the original package ID from the map.
	delete(ids, packageID(id))

	var results []source.PackageHandle
	for id := range ids {
		ph, err := s.packageHandle(ctx, id, source.ParseFull)
		if err != nil {
			return nil, err
		}
		results = append(results, ph)
	}
	return results, nil
}

// transitiveReverseDependencies populates the uris map with file URIs
// belonging to the provided package and its transitive reverse dependencies.
func (s *snapshot) transitiveReverseDependencies(id packageID, ids map[packageID]struct{}) {
	if _, ok := ids[id]; ok {
		return
	}
	if s.getMetadata(id) == nil {
		return
	}
	ids[id] = struct{}{}
	importedBy := s.getImportedBy(id)
	for _, parentID := range importedBy {
		s.transitiveReverseDependencies(parentID, ids)
	}
}

func (s *snapshot) getImportedBy(id packageID) []packageID {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getImportedByLocked(id)
}

func (s *snapshot) getImportedByLocked(id packageID) []packageID {
	// If we haven't rebuilt the import graph since creating the snapshot.
	if len(s.importedBy) == 0 {
		s.rebuildImportGraph()
	}
	return s.importedBy[id]
}

func (s *snapshot) clearAndRebuildImportGraph() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Completely invalidate the original map.
	s.importedBy = make(map[packageID][]packageID)
	s.rebuildImportGraph()
}

func (s *snapshot) rebuildImportGraph() {
	for id, m := range s.metadata {
		for _, importID := range m.deps {
			s.importedBy[importID] = append(s.importedBy[importID], id)
		}
	}
}

func (s *snapshot) addPackage(ph *packageHandle) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO: We should make sure not to compute duplicate packageHandles,
	// and instead panic here. This will be hard to do because we may encounter
	// the same package multiple times in the dependency tree.
	if _, ok := s.packages[ph.packageKey()]; ok {
		return
	}
	s.packages[ph.packageKey()] = ph
}

func (s *snapshot) workspacePackageIDs() (ids []packageID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id := range s.workspacePackages {
		ids = append(ids, id)
	}
	return ids
}

func (s *snapshot) WorkspacePackages(ctx context.Context) ([]source.PackageHandle, error) {
	if err := s.awaitLoaded(ctx); err != nil {
		return nil, err
	}
	var results []source.PackageHandle
	for _, pkgID := range s.workspacePackageIDs() {
		ph, err := s.packageHandle(ctx, pkgID, source.ParseFull)
		if err != nil {
			return nil, err
		}
		results = append(results, ph)
	}
	return results, nil
}

func (s *snapshot) KnownPackages(ctx context.Context) ([]source.PackageHandle, error) {
	if err := s.awaitLoaded(ctx); err != nil {
		return nil, err
	}
	// Collect PackageHandles for all of the workspace packages first.
	// They may need to be reloaded if their metadata has been invalidated.
	wsPackages := make(map[packageID]bool)
	s.mu.Lock()
	for id := range s.workspacePackages {
		wsPackages[id] = true
	}
	s.mu.Unlock()

	var results []source.PackageHandle
	for pkgID := range wsPackages {
		ph, err := s.packageHandle(ctx, pkgID, source.ParseFull)
		if err != nil {
			return nil, err
		}
		results = append(results, ph)
	}

	// Once all workspace packages have been checked, the metadata will be up-to-date.
	// Add all packages known in the workspace (that haven't already been added).
	pkgIDs := make(map[packageID]bool)
	s.mu.Lock()
	for id := range s.metadata {
		if !wsPackages[id] {
			pkgIDs[id] = true
		}
	}
	s.mu.Unlock()

	for pkgID := range pkgIDs {
		// Metadata for these packages should already be up-to-date,
		// so just build the package handle directly (without a reload).
		ph, err := s.buildPackageHandle(ctx, pkgID, source.ParseExported)
		if err != nil {
			return nil, err
		}
		results = append(results, ph)
	}
	return results, nil
}

func (s *snapshot) CachedImportPaths(ctx context.Context) (map[string]source.Package, error) {
	// Don't reload workspace package metadata.
	// This function is meant to only return currently cached information.
	s.view.awaitInitialized(ctx)

	s.mu.Lock()
	defer s.mu.Unlock()

	results := map[string]source.Package{}
	for _, ph := range s.packages {
		cachedPkg, err := ph.cached()
		if err != nil {
			continue
		}
		for importPath, newPkg := range cachedPkg.imports {
			if oldPkg, ok := results[string(importPath)]; ok {
				// Using the same trick as NarrowestPackageHandle, prefer non-variants.
				if len(newPkg.compiledGoFiles) < len(oldPkg.(*pkg).compiledGoFiles) {
					results[string(importPath)] = newPkg
				}
			} else {
				results[string(importPath)] = newPkg
			}
		}
	}
	return results, nil
}

func (s *snapshot) getPackage(id packageID, mode source.ParseMode) *packageHandle {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := packageKey{
		id:   id,
		mode: mode,
	}
	return s.packages[key]
}

func (s *snapshot) getActionHandle(id packageID, m source.ParseMode, a *analysis.Analyzer) *actionHandle {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := actionKey{
		pkg: packageKey{
			id:   id,
			mode: m,
		},
		analyzer: a,
	}
	return s.actions[key]
}

func (s *snapshot) addActionHandle(ah *actionHandle) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := actionKey{
		analyzer: ah.analyzer,
		pkg: packageKey{
			id:   ah.pkg.id,
			mode: ah.pkg.mode,
		},
	}
	if _, ok := s.actions[key]; ok {
		return
	}
	s.actions[key] = ah
}

func (s *snapshot) getIDsForURI(uri span.URI) []packageID {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ids[uri]
}

func (s *snapshot) getMetadataForURILocked(uri span.URI) (metadata []*metadata) {
	// TODO(matloob): uri can be a file or directory. Should we update the mappings
	// to map directories to their contained packages?

	for _, id := range s.ids[uri] {
		if m, ok := s.metadata[id]; ok {
			metadata = append(metadata, m)
		}
	}
	return metadata
}

func (s *snapshot) getMetadata(id packageID) *metadata {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.metadata[id]
}

func (s *snapshot) addID(uri span.URI, id packageID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existingID := range s.ids[uri] {
		if existingID == id {
			// TODO: We should make sure not to set duplicate IDs,
			// and instead panic here. This can be done by making sure not to
			// reset metadata information for packages we've already seen.
			return
		}
	}
	s.ids[uri] = append(s.ids[uri], id)
}

func (s *snapshot) isWorkspacePackage(id packageID) (packagePath, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	scope, ok := s.workspacePackages[id]
	return scope, ok
}

// GetFile returns a File for the given URI. It will always succeed because it
// adds the file to the managed set if needed.
func (s *snapshot) GetFile(uri span.URI) (source.FileHandle, error) {
	f, err := s.view.getFile(uri)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.files[f.URI()]; !ok {
		s.files[f.URI()] = s.view.session.cache.GetFile(uri)
	}
	return s.files[f.URI()], nil
}

func (s *snapshot) IsOpen(uri span.URI) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, open := s.files[uri].(*overlay)
	return open
}

func (s *snapshot) IsSaved(uri span.URI) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	ovl, open := s.files[uri].(*overlay)
	return !open || ovl.saved
}

func (s *snapshot) awaitLoaded(ctx context.Context) error {
	// Do not return results until the snapshot's view has been initialized.
	s.view.awaitInitialized(ctx)

	if err := s.reloadWorkspace(ctx); err != nil {
		return err
	}
	return s.reloadOrphanedFiles(ctx)
}

// reloadWorkspace reloads the metadata for all invalidated workspace packages.
func (s *snapshot) reloadWorkspace(ctx context.Context) error {
	// If the view's build configuration is invalid, we cannot reload by package path.
	// Just reload the directory instead.
	if !s.view.hasValidBuildConfiguration {
		return s.load(ctx, viewLoadScope("LOAD_INVALID_VIEW"))
	}

	// See which of the workspace packages are missing metadata.
	s.mu.Lock()
	var pkgPaths []interface{}
	for id, pkgPath := range s.workspacePackages {
		if s.metadata[id] == nil {
			pkgPaths = append(pkgPaths, pkgPath)
		}
	}
	s.mu.Unlock()

	if len(pkgPaths) == 0 {
		return nil
	}
	return s.load(ctx, pkgPaths...)
}

func (s *snapshot) reloadOrphanedFiles(ctx context.Context) error {
	// When we load ./... or a package path directly, we may not get packages
	// that exist only in overlays. As a workaround, we search all of the files
	// available in the snapshot and reload their metadata individually using a
	// file= query if the metadata is unavailable.
	scopes := s.orphanedFileScopes()
	if len(scopes) == 0 {
		return nil
	}

	err := s.load(ctx, scopes...)

	// If we failed to load some files, i.e. they have no metadata,
	// mark the failures so we don't bother retrying until the file's
	// content changes.
	//
	// TODO(rstambler): This may be an overestimate if the load stopped
	// early for an unrelated errors. Add a fallback?
	//
	// Check for context cancellation so that we don't incorrectly mark files
	// as unloadable, but don't return before setting all workspace packages.
	if ctx.Err() == nil && err != nil {
		s.mu.Lock()
		for _, scope := range scopes {
			uri := span.URI(scope.(fileURI))
			if s.getMetadataForURILocked(uri) == nil {
				s.unloadableFiles[uri] = struct{}{}
			}
		}
		s.mu.Unlock()
	}
	return nil
}

func (s *snapshot) orphanedFileScopes() []interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	scopeSet := make(map[span.URI]struct{})
	for uri, fh := range s.files {
		// Don't try to reload metadata for go.mod files.
		if fh.Identity().Kind != source.Go {
			continue
		}
		// If the URI doesn't belong to this view, then it's not in a workspace
		// package and should not be reloaded directly.
		if !contains(s.view.session.viewsOf(uri), s.view) {
			continue
		}
		// Don't reload metadata for files we've already deemed unloadable.
		if _, ok := s.unloadableFiles[uri]; ok {
			continue
		}
		if s.getMetadataForURILocked(uri) == nil {
			scopeSet[uri] = struct{}{}
		}
	}
	var scopes []interface{}
	for uri := range scopeSet {
		scopes = append(scopes, fileURI(uri))
	}
	return scopes
}

func contains(views []*view, view *view) bool {
	for _, v := range views {
		if v == view {
			return true
		}
	}
	return false
}

func (s *snapshot) clone(ctx context.Context, withoutURIs map[span.URI]source.FileHandle) *snapshot {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := &snapshot{
		id:                s.id + 1,
		view:              s.view,
		ids:               make(map[span.URI][]packageID),
		importedBy:        make(map[packageID][]packageID),
		metadata:          make(map[packageID]*metadata),
		packages:          make(map[packageKey]*packageHandle),
		actions:           make(map[actionKey]*actionHandle),
		files:             make(map[span.URI]source.FileHandle),
		workspacePackages: make(map[packageID]packagePath),
		unloadableFiles:   make(map[span.URI]struct{}),
	}

	// Copy all of the FileHandles.
	for k, v := range s.files {
		result.files[k] = v
	}
	// Copy the set of unloadable files.
	for k, v := range s.unloadableFiles {
		result.unloadableFiles[k] = v
	}

	// transitiveIDs keeps track of transitive reverse dependencies.
	// If an ID is present in the map, invalidate its types.
	// If an ID's value is true, invalidate its metadata too.
	transitiveIDs := make(map[packageID]bool)

	for withoutURI, currentFH := range withoutURIs {
		directIDs := map[packageID]struct{}{}

		// Collect all of the package IDs that correspond to the given file.
		// TODO: if the file has moved into a new package, we should invalidate that too.
		for _, id := range s.ids[withoutURI] {
			directIDs[id] = struct{}{}
		}
		// The original FileHandle for this URI is cached on the snapshot.
		originalFH := s.files[withoutURI]

		// Check if the file's package name or imports have changed,
		// and if so, invalidate this file's packages' metadata.
		invalidateMetadata := s.shouldInvalidateMetadata(ctx, originalFH, currentFH)

		// If a go.mod file's contents have changed, invalidate the metadata
		// for all of the packages in the workspace.
		if invalidateMetadata && currentFH.Identity().Kind == source.Mod {
			for id := range s.workspacePackages {
				directIDs[id] = struct{}{}
			}
		}

		// If this is a file we don't yet know about,
		// then we do not yet know what packages it should belong to.
		// Make a rough estimate of what metadata to invalidate by finding the package IDs
		// of all of the files in the same directory as this one.
		// TODO(rstambler): Speed this up by mapping directories to filenames.
		if len(directIDs) == 0 {
			if dirStat, err := os.Stat(filepath.Dir(withoutURI.Filename())); err == nil {
				for uri := range s.files {
					if fdirStat, err := os.Stat(filepath.Dir(uri.Filename())); err == nil {
						if os.SameFile(dirStat, fdirStat) {
							for _, id := range s.ids[uri] {
								directIDs[id] = struct{}{}
							}
						}
					}
				}
			}
		}

		// Invalidate reverse dependencies too.
		// TODO(heschi): figure out the locking model and use transitiveReverseDeps?
		var addRevDeps func(packageID)
		addRevDeps = func(id packageID) {
			if _, seen := transitiveIDs[id]; seen {
				return
			}
			transitiveIDs[id] = invalidateMetadata
			for _, rid := range s.getImportedByLocked(id) {
				addRevDeps(rid)
			}
		}
		for id := range directIDs {
			addRevDeps(id)
		}

		// Handle the invalidated file; it may have new contents or not exist.
		if _, _, err := currentFH.Read(ctx); os.IsNotExist(err) {
			delete(result.files, withoutURI)
		} else {
			result.files[withoutURI] = currentFH
		}
		// Make sure to remove the changed file from the unloadable set.
		delete(result.unloadableFiles, withoutURI)
	}

	// Collect the IDs for the packages associated with the excluded URIs.
	for k, ids := range s.ids {
		result.ids[k] = ids
	}
	// Copy the set of initally loaded packages.
	for k, v := range s.workspacePackages {
		result.workspacePackages[k] = v
	}
	// Copy the package type information.
	for k, v := range s.packages {
		if _, ok := transitiveIDs[k.id]; ok {
			continue
		}
		result.packages[k] = v
	}
	// Copy the package analysis information.
	for k, v := range s.actions {
		if _, ok := transitiveIDs[k.pkg.id]; ok {
			continue
		}
		result.actions[k] = v
	}
	// Copy the package metadata. We only need to invalidate packages directly
	// containing the affected file, and only if it changed in a relevant way.
	for k, v := range s.metadata {
		if invalidateMetadata, ok := transitiveIDs[k]; invalidateMetadata && ok {
			continue
		}
		result.metadata[k] = v
	}
	// Don't bother copying the importedBy graph,
	// as it changes each time we update metadata.

	return result
}

// shouldInvalidateMetadata reparses a file's package and import declarations to
// determine if the file requires a metadata reload.
func (s *snapshot) shouldInvalidateMetadata(ctx context.Context, originalFH, currentFH source.FileHandle) bool {
	if originalFH == nil {
		return currentFH.Identity().Kind == source.Go
	}
	// If the file hasn't changed, there's no need to reload.
	if originalFH.Identity().String() == currentFH.Identity().String() {
		return false
	}
	// If a go.mod file's contents have changed, always invalidate metadata.
	if kind := originalFH.Identity().Kind; kind == source.Mod {
		modfile, _ := s.view.ModFiles()
		return originalFH.Identity().URI == modfile
	}
	// Get the original and current parsed files in order to check package name and imports.
	original, _, _, _, originalErr := s.view.session.cache.ParseGoHandle(originalFH, source.ParseHeader).Parse(ctx)
	current, _, _, _, currentErr := s.view.session.cache.ParseGoHandle(currentFH, source.ParseHeader).Parse(ctx)
	if originalErr != nil || currentErr != nil {
		return (originalErr == nil) != (currentErr == nil)
	}
	// Check if the package's metadata has changed. The cases handled are:
	//    1. A package's name has changed
	//    2. A file's imports have changed
	if original.Name.Name != current.Name.Name {
		return true
	}
	// If the package's imports have increased, definitely re-run `go list`.
	if len(original.Imports) < len(current.Imports) {
		return true
	}
	importSet := make(map[string]struct{})
	for _, importSpec := range original.Imports {
		importSet[importSpec.Path.Value] = struct{}{}
	}
	// If any of the current imports were not in the original imports.
	for _, importSpec := range current.Imports {
		if _, ok := importSet[importSpec.Path.Value]; !ok {
			return true
		}
	}
	return false
}
