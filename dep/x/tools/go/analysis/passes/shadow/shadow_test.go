package shadow_test

import (
	"testing"

	"utilware/dep/x/tools/go/analysis/analysistest"
	"utilware/dep/x/tools/go/analysis/passes/shadow"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, shadow.Analyzer, "a")
}
