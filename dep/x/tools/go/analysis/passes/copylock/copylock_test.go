// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package copylock_test

import (
	"testing"

	"utilware/dep/x/tools/go/analysis/analysistest"
	"utilware/dep/x/tools/go/analysis/passes/copylock"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, copylock.Analyzer, "a")
}
