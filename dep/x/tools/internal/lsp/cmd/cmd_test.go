// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"

	"utilware/dep/x/tools/go/packages/packagestest"
	"utilware/dep/x/tools/internal/jsonrpc2/servertest"
	"utilware/dep/x/tools/internal/lsp/cache"
	"utilware/dep/x/tools/internal/lsp/cmd"
	cmdtest "utilware/dep/x/tools/internal/lsp/cmd/test"
	"utilware/dep/x/tools/internal/lsp/lsprpc"
	"utilware/dep/x/tools/internal/lsp/tests"
	"utilware/dep/x/tools/internal/testenv"
)

func TestMain(m *testing.M) {
	testenv.ExitIfSmallMachine()
	os.Exit(m.Run())
}

func TestCommandLine(t *testing.T) {
	packagestest.TestAll(t, testCommandLine)
}

func testCommandLine(t *testing.T, exporter packagestest.Exporter) {
	ctx := tests.Context(t)
	ts := testServer(ctx)
	data := tests.Load(t, exporter, "../testdata")
	for _, datum := range data {
		defer datum.Exported.Cleanup()
		t.Run(datum.Folder, func(t *testing.T) {
			t.Helper()
			tests.Run(t, cmdtest.NewRunner(exporter, datum, ctx, ts.Addr, nil), datum)
		})
	}
}

func testServer(ctx context.Context) *servertest.TCPServer {
	cache := cache.New(nil)
	ss := lsprpc.NewStreamServer(cache, false)
	return servertest.NewTCPServer(ctx, ss)
}

func TestDefinitionHelpExample(t *testing.T) {
	// TODO: https://ccms/issue/32794.
	t.Skip()
	if runtime.GOOS == "android" {
		t.Skip("not all source files are available on android")
	}
	dir, err := os.Getwd()
	if err != nil {
		t.Errorf("could not get wd: %v", err)
		return
	}
	ctx := tests.Context(t)
	ts := testServer(ctx)
	thisFile := filepath.Join(dir, "definition.go")
	baseArgs := []string{"query", "definition"}
	expect := regexp.MustCompile(`(?s)^[\w/\\:_-]+flag[/\\]flag.go:\d+:\d+-\d+: defined here as FlagSet struct {.*}$`)
	for _, query := range []string{
		fmt.Sprintf("%v:%v:%v", thisFile, cmd.ExampleLine, cmd.ExampleColumn),
		fmt.Sprintf("%v:#%v", thisFile, cmd.ExampleOffset)} {
		args := append(baseArgs, query)
		r := cmdtest.NewRunner(nil, nil, ctx, ts.Addr, nil)
		got, _ := r.NormalizeGoplsCmd(t, args...)
		if !expect.MatchString(got) {
			t.Errorf("test with %v\nexpected:\n%s\ngot:\n%s", args, expect, got)
		}
	}
}
