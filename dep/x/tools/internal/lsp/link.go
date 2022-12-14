// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"utilware/dep/x/tools/internal/lsp/protocol"
	"utilware/dep/x/tools/internal/lsp/source"
	"utilware/dep/x/tools/internal/span"
	"utilware/dep/x/tools/internal/telemetry/log"
)

func (s *Server) documentLink(ctx context.Context, params *protocol.DocumentLinkParams) ([]protocol.DocumentLink, error) {
	// TODO(ccms/go#36501): Support document links for go.mod files.
	snapshot, fh, ok, err := s.beginFileRequest(params.TextDocument.URI, source.Go)
	if !ok {
		return nil, err
	}
	view := snapshot.View()
	phs, err := view.Snapshot().PackageHandles(ctx, fh)
	if err != nil {
		return nil, err
	}
	ph, err := source.WidestPackageHandle(phs)
	if err != nil {
		return nil, err
	}
	file, _, m, _, err := view.Session().Cache().ParseGoHandle(fh, source.ParseFull).Parse(ctx)
	if err != nil {
		return nil, err
	}
	var links []protocol.DocumentLink
	ast.Inspect(file, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.ImportSpec:
			// For import specs, provide a link to a documentation website, like https://pkg.go.dev.
			if target, err := strconv.Unquote(n.Path.Value); err == nil {
				if mod, version, ok := moduleAtVersion(ctx, target, ph); ok && strings.ToLower(view.Options().LinkTarget) == "pkg.go.dev" {
					target = strings.Replace(target, mod, mod+"@"+version, 1)
				}
				target = fmt.Sprintf("https://%s/%s", view.Options().LinkTarget, target)
				// Account for the quotation marks in the positions.
				start, end := n.Path.Pos()+1, n.Path.End()-1
				if l, err := toProtocolLink(view, m, target, start, end); err == nil {
					links = append(links, l)
				} else {
					log.Error(ctx, "failed to create protocol link", err)
				}
			}
			return false
		case *ast.BasicLit:
			// Look for links in string literals.
			if n.Kind == token.STRING {
				links = append(links, findLinksInString(ctx, view, n.Value, n.Pos(), m)...)
			}
			return false
		}
		return true
	})
	// Look for links in comments.
	for _, commentGroup := range file.Comments {
		for _, comment := range commentGroup.List {
			links = append(links, findLinksInString(ctx, view, comment.Text, comment.Pos(), m)...)
		}
	}
	return links, nil
}

func moduleAtVersion(ctx context.Context, target string, ph source.PackageHandle) (string, string, bool) {
	pkg, err := ph.Check(ctx)
	if err != nil {
		return "", "", false
	}
	impPkg, err := pkg.GetImport(target)
	if err != nil {
		return "", "", false
	}
	if impPkg.Module() == nil {
		return "", "", false
	}
	version, modpath := impPkg.Module().Version, impPkg.Module().Path
	if modpath == "" || version == "" {
		return "", "", false
	}
	return modpath, version, true
}

func findLinksInString(ctx context.Context, view source.View, src string, pos token.Pos, m *protocol.ColumnMapper) []protocol.DocumentLink {
	var links []protocol.DocumentLink
	for _, index := range view.Options().URLRegexp.FindAllIndex([]byte(src), -1) {
		start, end := index[0], index[1]
		startPos := token.Pos(int(pos) + start)
		endPos := token.Pos(int(pos) + end)
		url, err := url.Parse(src[start:end])
		if err != nil {
			log.Error(ctx, "failed to parse matching URL", err)
			continue
		}
		// If the URL has no scheme, use https.
		if url.Scheme == "" {
			url.Scheme = "https"
		}
		l, err := toProtocolLink(view, m, url.String(), startPos, endPos)
		if err != nil {
			log.Error(ctx, "failed to create protocol link", err)
			continue
		}
		links = append(links, l)
	}
	// Handle ccms/go#1234-style links.
	r := getIssueRegexp()
	for _, index := range r.FindAllIndex([]byte(src), -1) {
		start, end := index[0], index[1]
		startPos := token.Pos(int(pos) + start)
		endPos := token.Pos(int(pos) + end)
		matches := r.FindStringSubmatch(src)
		if len(matches) < 4 {
			continue
		}
		org, repo, number := matches[1], matches[2], matches[3]
		target := fmt.Sprintf("https://github.com/%s/%s/issues/%s", org, repo, number)
		l, err := toProtocolLink(view, m, target, startPos, endPos)
		if err != nil {
			log.Error(ctx, "failed to create protocol link", err)
			continue
		}
		links = append(links, l)
	}
	return links
}

func getIssueRegexp() *regexp.Regexp {
	once.Do(func() {
		issueRegexp = regexp.MustCompile(`(\w+)/([\w-]+)#([0-9]+)`)
	})
	return issueRegexp
}

var (
	once        sync.Once
	issueRegexp *regexp.Regexp
)

func toProtocolLink(view source.View, m *protocol.ColumnMapper, target string, start, end token.Pos) (protocol.DocumentLink, error) {
	spn, err := span.NewRange(view.Session().Cache().FileSet(), start, end).Span()
	if err != nil {
		return protocol.DocumentLink{}, err
	}
	rng, err := m.Range(spn)
	if err != nil {
		return protocol.DocumentLink{}, err
	}
	return protocol.DocumentLink{
		Range:  rng,
		Target: target,
	}, nil
}
