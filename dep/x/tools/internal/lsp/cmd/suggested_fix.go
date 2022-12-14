// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"

	"utilware/dep/x/tools/internal/lsp/diff"
	"utilware/dep/x/tools/internal/lsp/protocol"
	"utilware/dep/x/tools/internal/lsp/source"
	"utilware/dep/x/tools/internal/span"
	"utilware/dep/x/tools/internal/tool"
	errors "utilware/dep/x/xerrors"
)

// suggestedfix implements the fix verb for gopls.
type suggestedfix struct {
	Diff  bool `flag:"d" help:"display diffs instead of rewriting files"`
	Write bool `flag:"w" help:"write result to (source) file instead of stdout"`
	All   bool `flag:"a" help:"apply all fixes, not just preferred fixes"`

	app *Application
}

func (s *suggestedfix) Name() string      { return "fix" }
func (s *suggestedfix) Usage() string     { return "<filename>" }
func (s *suggestedfix) ShortHelp() string { return "apply suggested fixes" }
func (s *suggestedfix) DetailedHelp(f *flag.FlagSet) {
	fmt.Fprintf(f.Output(), `
Example: apply suggested fixes for this file:

  $ gopls fix -w internal/lsp/cmd/check.go

gopls fix flags are:
`)
	f.PrintDefaults()
}

// Run performs diagnostic checks on the file specified and either;
// - if -w is specified, updates the file in place;
// - if -d is specified, prints out unified diffs of the changes; or
// - otherwise, prints the new versions to stdout.
func (s *suggestedfix) Run(ctx context.Context, args ...string) error {
	if len(args) != 1 {
		return tool.CommandLineErrorf("fix expects 1 argument")
	}
	conn, err := s.app.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.terminate(ctx)

	from := span.Parse(args[0])
	uri := from.URI()
	file := conn.AddFile(ctx, uri)
	if file.err != nil {
		return file.err
	}

	if err := conn.diagnoseFiles(ctx, []span.URI{uri}); err != nil {
		return err
	}
	conn.Client.filesMu.Lock()
	defer conn.Client.filesMu.Unlock()

	p := protocol.CodeActionParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: protocol.URIFromSpanURI(uri),
		},
		Context: protocol.CodeActionContext{
			Only:        []protocol.CodeActionKind{protocol.QuickFix},
			Diagnostics: file.diagnostics,
		},
	}
	actions, err := conn.CodeAction(ctx, &p)
	if err != nil {
		return errors.Errorf("%v: %v", from, err)
	}
	var edits []protocol.TextEdit
	for _, a := range actions {
		if !a.IsPreferred && !s.All {
			continue
		}
		for _, c := range a.Edit.DocumentChanges {
			if fileURI(c.TextDocument.URI) == uri {
				edits = append(edits, c.Edits...)
			}
		}
	}

	sedits, err := source.FromProtocolEdits(file.mapper, edits)
	if err != nil {
		return errors.Errorf("%v: %v", edits, err)
	}
	newContent := diff.ApplyEdits(string(file.mapper.Content), sedits)

	filename := file.uri.Filename()
	switch {
	case s.Write:
		if len(edits) > 0 {
			ioutil.WriteFile(filename, []byte(newContent), 0644)
		}
	case s.Diff:
		diffs := diff.ToUnified(filename+".orig", filename, string(file.mapper.Content), sedits)
		fmt.Print(diffs)
	default:
		fmt.Print(string(newContent))
	}
	return nil
}
