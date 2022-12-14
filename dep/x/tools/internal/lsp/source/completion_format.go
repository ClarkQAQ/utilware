// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package source

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/printer"
	"go/types"
	"strings"

	"utilware/dep/x/tools/internal/imports"
	"utilware/dep/x/tools/internal/lsp/protocol"
	"utilware/dep/x/tools/internal/lsp/snippet"
	"utilware/dep/x/tools/internal/lsp/telemetry"
	"utilware/dep/x/tools/internal/span"
	"utilware/dep/x/tools/internal/telemetry/log"
	"utilware/dep/x/tools/internal/telemetry/tag"
	errors "utilware/dep/x/xerrors"
)

// formatCompletion creates a completion item for a given candidate.
func (c *completer) item(cand candidate) (CompletionItem, error) {
	obj := cand.obj

	// Handle builtin types separately.
	if obj.Parent() == types.Universe {
		return c.formatBuiltin(cand), nil
	}

	var (
		label         = cand.name
		detail        = types.TypeString(obj.Type(), c.qf)
		insert        = label
		kind          = protocol.TextCompletion
		snip          *snippet.Builder
		protocolEdits []protocol.TextEdit
	)
	if obj.Type() == nil {
		detail = ""
	}

	// expandFuncCall mutates the completion label, detail, and snippet
	// to that of an invocation of sig.
	expandFuncCall := func(sig *types.Signature) {
		params := formatParams(c.ctx, c.snapshot, c.pkg, sig, c.qf)
		snip = c.functionCallSnippet(label, params)
		results, writeParens := formatResults(sig.Results(), c.qf)
		detail = "func" + formatFunction(params, results, writeParens)

		// Add variadic "..." if we are using a function result to fill in a variadic parameter.
		if sig.Results().Len() == 1 && c.inference.matchesVariadic(sig.Results().At(0).Type()) {
			snip.WriteText("...")
		}
	}

	switch obj := obj.(type) {
	case *types.TypeName:
		detail, kind = formatType(obj.Type(), c.qf)
	case *types.Const:
		kind = protocol.ConstantCompletion
	case *types.Var:
		if _, ok := obj.Type().(*types.Struct); ok {
			detail = "struct{...}" // for anonymous structs
		} else if obj.IsField() {
			var err error
			detail, err = formatFieldType(c.ctx, c.snapshot, c.pkg, obj, c.qf)
			if err != nil {
				detail = types.TypeString(obj.Type(), c.qf)
			}
		}
		if obj.IsField() {
			kind = protocol.FieldCompletion
			snip = c.structFieldSnippet(label, detail)
		} else {
			kind = protocol.VariableCompletion
		}
		if obj.Type() == nil {
			break
		}

		if sig, ok := obj.Type().Underlying().(*types.Signature); ok && cand.expandFuncCall {
			expandFuncCall(sig)
		}

		// Add variadic "..." if we are using a variable to fill in a variadic parameter.
		if c.inference.matchesVariadic(obj.Type()) {
			snip = &snippet.Builder{}
			snip.WriteText(insert + "...")
		}
	case *types.Func:
		sig, ok := obj.Type().Underlying().(*types.Signature)
		if !ok {
			break
		}
		kind = protocol.FunctionCompletion
		if sig != nil && sig.Recv() != nil {
			kind = protocol.MethodCompletion
		}

		if cand.expandFuncCall {
			expandFuncCall(sig)
		}
	case *types.PkgName:
		kind = protocol.ModuleCompletion
		detail = fmt.Sprintf("%q", obj.Imported().Path())
	case *types.Label:
		kind = protocol.ConstantCompletion
		detail = "label"
	}

	// If this candidate needs an additional import statement,
	// add the additional text edits needed.
	if cand.imp != nil {
		addlEdits, err := c.importEdits(cand.imp)
		if err != nil {
			return CompletionItem{}, err
		}

		protocolEdits = append(protocolEdits, addlEdits...)
		if kind != protocol.ModuleCompletion {
			if detail != "" {
				detail += " "
			}
			detail += fmt.Sprintf("(from %q)", cand.imp.importPath)
		}
	}

	// Prepend "&" or "*" operator as appropriate.
	var prefixOp string
	if cand.takeAddress {
		prefixOp = "&"
	} else if cand.makePointer {
		prefixOp = "*"
	} else if cand.dereference > 0 {
		prefixOp = strings.Repeat("*", cand.dereference)
	}

	if prefixOp != "" {
		// If we are in a selector, add an edit to place prefix before selector.
		if sel := enclosingSelector(c.path, c.pos); sel != nil {
			edits, err := prependEdit(c.snapshot.View().Session().Cache().FileSet(), c.mapper, sel, prefixOp)
			if err != nil {
				log.Error(c.ctx, "error generating prefix edit", err)
			} else {
				protocolEdits = append(protocolEdits, edits...)
			}
		} else {
			// If there is no selector, just stick the prefix at the start.
			insert = prefixOp + insert
		}

		label = prefixOp + label
	}

	detail = strings.TrimPrefix(detail, "untyped ")
	item := CompletionItem{
		Label:               label,
		InsertText:          insert,
		AdditionalTextEdits: protocolEdits,
		Detail:              detail,
		Kind:                kind,
		Score:               cand.score,
		Depth:               len(c.deepState.chain),
		snippet:             snip,
	}
	// If the user doesn't want documentation for completion items.
	if !c.opts.documentation {
		return item, nil
	}
	pos := c.snapshot.View().Session().Cache().FileSet().Position(obj.Pos())

	// We ignore errors here, because some types, like "unsafe" or "error",
	// may not have valid positions that we can use to get documentation.
	if !pos.IsValid() {
		return item, nil
	}
	uri := span.URIFromPath(pos.Filename)

	// Find the source file of the candidate, starting from a package
	// that should have it in its dependencies.
	searchPkg := c.pkg
	if cand.imp != nil && cand.imp.pkg != nil {
		searchPkg = cand.imp.pkg
	}
	file, pkg, err := findPosInPackage(c.snapshot.View(), searchPkg, obj.Pos())
	if err != nil {
		return item, nil
	}
	ident, err := findIdentifier(c.ctx, c.snapshot, pkg, file, obj.Pos())
	if err != nil {
		return item, nil
	}
	hover, err := ident.Hover(c.ctx)
	if err != nil {
		log.Error(c.ctx, "failed to find Hover", err, telemetry.URI.Of(uri))
		return item, nil
	}
	item.Documentation = hover.Synopsis
	if c.opts.fullDocumentation {
		item.Documentation = hover.FullDocumentation
	}
	return item, nil
}

// importEdits produces the text edits necessary to add the given import to the current file.
func (c *completer) importEdits(imp *importInfo) ([]protocol.TextEdit, error) {
	if imp == nil {
		return nil, nil
	}

	uri := span.URIFromPath(c.filename)
	var ph ParseGoHandle
	for _, h := range c.pkg.CompiledGoFiles() {
		if h.File().Identity().URI == uri {
			ph = h
		}
	}
	if ph == nil {
		return nil, errors.Errorf("building import completion for %v: no ParseGoHandle for %s", imp.importPath, c.filename)
	}

	return computeOneImportFixEdits(c.ctx, c.snapshot.View(), ph, &imports.ImportFix{
		StmtInfo: imports.ImportInfo{
			ImportPath: imp.importPath,
			Name:       imp.name,
		},
		// IdentName is unused on this path and is difficult to get.
		FixType: imports.AddImport,
	})
}

func (c *completer) formatBuiltin(cand candidate) CompletionItem {
	obj := cand.obj
	item := CompletionItem{
		Label:      obj.Name(),
		InsertText: obj.Name(),
		Score:      cand.score,
	}
	switch obj.(type) {
	case *types.Const:
		item.Kind = protocol.ConstantCompletion
	case *types.Builtin:
		item.Kind = protocol.FunctionCompletion
		astObj, err := c.snapshot.View().LookupBuiltin(c.ctx, obj.Name())
		if err != nil {
			if c.ctx.Err() == nil {
				log.Error(c.ctx, "no builtin package", err)
			}
			break
		}
		decl, ok := astObj.Decl.(*ast.FuncDecl)
		if !ok {
			break
		}
		params, _ := formatFieldList(c.ctx, c.snapshot.View(), decl.Type.Params)
		results, writeResultParens := formatFieldList(c.ctx, c.snapshot.View(), decl.Type.Results)
		item.Label = obj.Name()
		item.Detail = "func" + formatFunction(params, results, writeResultParens)
		item.snippet = c.functionCallSnippet(obj.Name(), params)
	case *types.TypeName:
		if types.IsInterface(obj.Type()) {
			item.Kind = protocol.InterfaceCompletion
		} else {
			item.Kind = protocol.ClassCompletion
		}
	case *types.Nil:
		item.Kind = protocol.VariableCompletion
	}
	return item
}

var replacer = strings.NewReplacer(
	`ComplexType`, `complex128`,
	`FloatType`, `float64`,
	`IntegerType`, `int`,
)

func formatFieldList(ctx context.Context, v View, list *ast.FieldList) ([]string, bool) {
	if list == nil {
		return nil, false
	}
	var writeResultParens bool
	var result []string
	for i := 0; i < len(list.List); i++ {
		if i >= 1 {
			writeResultParens = true
		}
		p := list.List[i]
		cfg := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 4}
		b := &bytes.Buffer{}
		if err := cfg.Fprint(b, v.Session().Cache().FileSet(), p.Type); err != nil {
			log.Error(ctx, "unable to print type", nil, tag.Of("Type", p.Type))
			continue
		}
		typ := replacer.Replace(b.String())
		if len(p.Names) == 0 {
			result = append(result, typ)
		}
		for _, name := range p.Names {
			if name.Name != "" {
				if i == 0 {
					writeResultParens = true
				}
				result = append(result, fmt.Sprintf("%s %s", name.Name, typ))
			} else {
				result = append(result, typ)
			}
		}
	}
	return result, writeResultParens
}

// qualifier returns a function that appropriately formats a types.PkgName
// appearing in a *ast.File.
func qualifier(f *ast.File, pkg *types.Package, info *types.Info) types.Qualifier {
	// Construct mapping of import paths to their defined or implicit names.
	imports := make(map[*types.Package]string)
	for _, imp := range f.Imports {
		var obj types.Object
		if imp.Name != nil {
			obj = info.Defs[imp.Name]
		} else {

			obj = info.Implicits[imp]
		}
		if pkgname, ok := obj.(*types.PkgName); ok {
			imports[pkgname.Imported()] = pkgname.Name()
		}
	}
	// Define qualifier to replace full package paths with names of the imports.
	return func(p *types.Package) string {
		if p == pkg {
			return ""
		}
		if name, ok := imports[p]; ok {
			return name
		}
		return p.Name()
	}
}
