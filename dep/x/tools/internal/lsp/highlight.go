// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

import (
	"context"

	"utilware/dep/x/tools/internal/lsp/protocol"
	"utilware/dep/x/tools/internal/lsp/source"
	"utilware/dep/x/tools/internal/lsp/telemetry"
	"utilware/dep/x/tools/internal/telemetry/log"
)

func (s *Server) documentHighlight(ctx context.Context, params *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	snapshot, fh, ok, err := s.beginFileRequest(params.TextDocument.URI, source.Go)
	if !ok {
		return nil, err
	}
	rngs, err := source.Highlight(ctx, snapshot, fh, params.Position)
	if err != nil {
		log.Error(ctx, "no highlight", err, telemetry.URI.Of(params.TextDocument.URI))
	}
	return toProtocolHighlight(rngs), nil
}

func toProtocolHighlight(rngs []protocol.Range) []protocol.DocumentHighlight {
	result := make([]protocol.DocumentHighlight, 0, len(rngs))
	kind := protocol.Text
	for _, rng := range rngs {
		result = append(result, protocol.DocumentHighlight{
			Kind:  kind,
			Range: rng,
		})
	}
	return result
}
