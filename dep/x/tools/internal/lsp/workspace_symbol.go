// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

import (
	"context"

	"utilware/dep/x/tools/internal/lsp/protocol"
	"utilware/dep/x/tools/internal/lsp/source"
	"utilware/dep/x/tools/internal/telemetry/trace"
)

func (s *Server) symbol(ctx context.Context, params *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) {
	ctx, done := trace.StartSpan(ctx, "lsp.Server.symbol")
	defer done()

	return source.WorkspaceSymbols(ctx, s.session.Views(), params.Query)
}
