// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cache

import (
	"context"
	"crypto/sha1"
	"fmt"
	"go/token"
	"reflect"
	"strconv"
	"sync/atomic"

	"utilware/dep/x/tools/internal/lsp/debug"
	"utilware/dep/x/tools/internal/lsp/source"
	"utilware/dep/x/tools/internal/memoize"
	"utilware/dep/x/tools/internal/span"
)

func New(options func(*source.Options)) *Cache {
	index := atomic.AddInt64(&cacheIndex, 1)
	c := &Cache{
		fs:      &nativeFileSystem{},
		id:      strconv.FormatInt(index, 10),
		fset:    token.NewFileSet(),
		options: options,
	}
	debug.AddCache(debugCache{c})
	return c
}

type Cache struct {
	fs      source.FileSystem
	id      string
	fset    *token.FileSet
	options func(*source.Options)

	store memoize.Store
}

type fileKey struct {
	identity source.FileIdentity
}

type fileHandle struct {
	cache      *Cache
	underlying source.FileHandle
	handle     *memoize.Handle
}

type fileData struct {
	memoize.NoCopy
	bytes []byte
	hash  string
	err   error
}

func (c *Cache) GetFile(uri span.URI) source.FileHandle {
	underlying := c.fs.GetFile(uri)
	key := fileKey{
		identity: underlying.Identity(),
	}
	h := c.store.Bind(key, func(ctx context.Context) interface{} {
		data := &fileData{}
		data.bytes, data.hash, data.err = underlying.Read(ctx)
		return data
	})
	return &fileHandle{
		cache:      c,
		underlying: underlying,
		handle:     h,
	}
}

func (c *Cache) NewSession() *Session {
	index := atomic.AddInt64(&sessionIndex, 1)
	s := &Session{
		cache:    c,
		id:       strconv.FormatInt(index, 10),
		options:  source.DefaultOptions(),
		overlays: make(map[span.URI]*overlay),
	}
	debug.AddSession(debugSession{s})
	return s
}

func (c *Cache) FileSet() *token.FileSet {
	return c.fset
}

func (h *fileHandle) FileSystem() source.FileSystem {
	return h.cache
}

func (h *fileHandle) Identity() source.FileIdentity {
	return h.underlying.Identity()
}

func (h *fileHandle) Read(ctx context.Context) ([]byte, string, error) {
	v := h.handle.Get(ctx)
	if v == nil {
		return nil, "", ctx.Err()
	}
	data := v.(*fileData)
	return data.bytes, data.hash, data.err
}

func hashContents(contents []byte) string {
	// TODO: consider whether sha1 is the best choice here
	// This hash is used for internal identity detection only
	return fmt.Sprintf("%x", sha1.Sum(contents))
}

var cacheIndex, sessionIndex, viewIndex int64

type debugCache struct{ *Cache }

func (c *Cache) ID() string                         { return c.id }
func (c debugCache) FileSet() *token.FileSet        { return c.fset }
func (c debugCache) MemStats() map[reflect.Type]int { return c.store.Stats() }
