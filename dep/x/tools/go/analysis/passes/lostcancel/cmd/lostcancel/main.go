// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The lostcancel command applies the utilware/dep/x/tools/go/analysis/passes/lostcancel
// analysis to the specified packages of Go source code.
package main

import (
	"utilware/dep/x/tools/go/analysis/passes/lostcancel"
	"utilware/dep/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(lostcancel.Analyzer) }
