// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.11

package packagestest_test

import (
	"path/filepath"
	"testing"

	"utilware/dep/x/tools/go/packages/packagestest"
)

func TestModulesExport(t *testing.T) {
	exported := packagestest.Export(t, packagestest.Modules, testdata)
	defer exported.Cleanup()
	// Check that the cfg contains all the right bits
	var expectDir = filepath.Join(exported.Temp(), "fake1")
	if exported.Config.Dir != expectDir {
		t.Errorf("Got working directory %v expected %v", exported.Config.Dir, expectDir)
	}
	checkFiles(t, exported, []fileTest{
		{"ccms/fake1", "go.mod", "fake1/go.mod", nil},
		{"ccms/fake1", "a.go", "fake1/a.go", checkLink("testdata/a.go")},
		{"ccms/fake1", "b.go", "fake1/b.go", checkContent("package fake1")},
		{"ccms/fake2", "go.mod", "modcache/pkg/mod/ccms/fake2@v1.0.0/go.mod", nil},
		{"ccms/fake2", "other/a.go", "modcache/pkg/mod/ccms/fake2@v1.0.0/other/a.go", checkContent("package fake2")},
		{"ccms/fake2/v2", "other/a.go", "modcache/pkg/mod/ccms/fake2/v2@v2.0.0/other/a.go", checkContent("package fake2")},
		{"ccms/fake3@v1.1.0", "other/a.go", "modcache/pkg/mod/ccms/fake3@v1.1.0/other/a.go", checkContent("package fake3")},
		{"ccms/fake3@v1.0.0", "other/a.go", "modcache/pkg/mod/ccms/fake3@v1.0.0/other/a.go", nil},
	})
}
