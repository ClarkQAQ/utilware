// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cache

import (
	"context"
	"io/ioutil"
	"os"

	"utilware/dep/x/tools/internal/lsp/source"
	"utilware/dep/x/tools/internal/lsp/telemetry"
	"utilware/dep/x/tools/internal/span"
	"utilware/dep/x/tools/internal/telemetry/trace"
	errors "utilware/dep/x/xerrors"
)

// ioLimit limits the number of parallel file reads per process.
var ioLimit = make(chan struct{}, 128)

// nativeFileSystem implements FileSystem reading from the normal os file system.
type nativeFileSystem struct{}

// nativeFileHandle implements FileHandle for nativeFileSystem
type nativeFileHandle struct {
	fs       *nativeFileSystem
	identity source.FileIdentity
}

func (fs *nativeFileSystem) GetFile(uri span.URI) source.FileHandle {
	return &nativeFileHandle{
		fs: fs,
		identity: source.FileIdentity{
			URI:        uri,
			Identifier: identifier(uri.Filename()),
			Kind:       source.DetectLanguage("", uri.Filename()),
		},
	}
}

func (h *nativeFileHandle) FileSystem() source.FileSystem {
	return h.fs
}

func (h *nativeFileHandle) Identity() source.FileIdentity {
	return h.identity
}

func (h *nativeFileHandle) Read(ctx context.Context) ([]byte, string, error) {
	ctx, done := trace.StartSpan(ctx, "cache.nativeFileHandle.Read", telemetry.File.Of(h.identity.URI.Filename()))
	_ = ctx
	defer done()

	ioLimit <- struct{}{}
	defer func() { <-ioLimit }()

	if id := identifier(h.identity.URI.Filename()); id != h.identity.Identifier {
		return nil, "", errors.Errorf("%s: file has been modified", h.identity.URI.Filename())
	}
	data, err := ioutil.ReadFile(h.identity.URI.Filename())
	if err != nil {
		return nil, "", err
	}
	return data, hashContents(data), nil
}

func identifier(filename string) string {
	if fi, err := os.Stat(filename); err == nil {
		return fi.ModTime().String()
	}
	return "DOES NOT EXIST"
}
