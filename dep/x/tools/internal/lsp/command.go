package lsp

import (
	"context"

	"utilware/dep/x/tools/internal/lsp/protocol"
	"utilware/dep/x/tools/internal/lsp/source"
	errors "utilware/dep/x/xerrors"
)

func (s *Server) executeCommand(ctx context.Context, params *protocol.ExecuteCommandParams) (interface{}, error) {
	switch params.Command {
	case "tidy":
		if len(params.Arguments) == 0 || len(params.Arguments) > 1 {
			return nil, errors.Errorf("expected one file URI for call to `go mod tidy`, got %v", params.Arguments)
		}
		uri := protocol.DocumentURI(params.Arguments[0].(string))
		snapshot, _, ok, err := s.beginFileRequest(uri, source.Mod)
		if !ok {
			return nil, err
		}
		// Run go.mod tidy on the view.
		if _, err := source.InvokeGo(ctx, snapshot.View().Folder().Filename(), snapshot.Config(ctx).Env, "mod", "tidy"); err != nil {
			return nil, err
		}
	case "upgrade.dependency":
		if len(params.Arguments) < 2 {
			return nil, errors.Errorf("expected one file URI and one dependency for call to `go get`, got %v", params.Arguments)
		}
		uri := protocol.DocumentURI(params.Arguments[0].(string))
		snapshot, _, ok, err := s.beginFileRequest(uri, source.UnknownKind)
		if !ok {
			return nil, err
		}
		dep := params.Arguments[1].(string)
		// Run "go get" on the dependency to upgrade it to the latest version.
		if _, err := source.InvokeGo(ctx, snapshot.View().Folder().Filename(), snapshot.Config(ctx).Env, "get", dep); err != nil {
			return nil, err
		}
	}
	return nil, nil
}
