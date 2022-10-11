// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

// This file provides an example command for static checkers
// conforming to the utilware/dep/x/tools/go/analysis API.
// It serves as a model for the behavior of the cmd/vet tool in $GOROOT.
// Being based on the unitchecker driver, it must be run by go vet:
//
//   $ go build -o unitchecker main.go
//   $ go vet -vettool=unitchecker my/project/...
//
// For a checker also capable of running standalone, use multichecker.
package main

import (
	"utilware/dep/x/tools/go/analysis/unitchecker"

	"utilware/dep/x/tools/go/analysis/passes/asmdecl"
	"utilware/dep/x/tools/go/analysis/passes/assign"
	"utilware/dep/x/tools/go/analysis/passes/atomic"
	"utilware/dep/x/tools/go/analysis/passes/bools"
	"utilware/dep/x/tools/go/analysis/passes/buildtag"
	"utilware/dep/x/tools/go/analysis/passes/cgocall"
	"utilware/dep/x/tools/go/analysis/passes/composite"
	"utilware/dep/x/tools/go/analysis/passes/copylock"
	"utilware/dep/x/tools/go/analysis/passes/errorsas"
	"utilware/dep/x/tools/go/analysis/passes/httpresponse"
	"utilware/dep/x/tools/go/analysis/passes/loopclosure"
	"utilware/dep/x/tools/go/analysis/passes/lostcancel"
	"utilware/dep/x/tools/go/analysis/passes/nilfunc"
	"utilware/dep/x/tools/go/analysis/passes/printf"
	"utilware/dep/x/tools/go/analysis/passes/shift"
	"utilware/dep/x/tools/go/analysis/passes/stdmethods"
	"utilware/dep/x/tools/go/analysis/passes/structtag"
	"utilware/dep/x/tools/go/analysis/passes/tests"
	"utilware/dep/x/tools/go/analysis/passes/unmarshal"
	"utilware/dep/x/tools/go/analysis/passes/unreachable"
	"utilware/dep/x/tools/go/analysis/passes/unsafeptr"
	"utilware/dep/x/tools/go/analysis/passes/unusedresult"
)

func main() {
	unitchecker.Main(
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		errorsas.Analyzer,
		httpresponse.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shift.Analyzer,
		stdmethods.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
	)
}
