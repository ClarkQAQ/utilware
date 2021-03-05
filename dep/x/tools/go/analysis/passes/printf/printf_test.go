package printf_test

import (
	"testing"

	"utilware/dep/x/tools/go/analysis/analysistest"
	"utilware/dep/x/tools/go/analysis/passes/printf"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	printf.Analyzer.Flags.Set("funcs", "Warn,Warnf")
	analysistest.Run(t, testdata, printf.Analyzer, "a", "b", "nofmt")
}
