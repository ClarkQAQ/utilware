// Copyright 2020 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !libc.dmesg

package libc // import "utilware/sqlite/libc"

const dmesgs = false

func dmesg(s string, args ...interface{}) {}
