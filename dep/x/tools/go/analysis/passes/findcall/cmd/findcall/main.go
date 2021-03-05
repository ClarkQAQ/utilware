// The findcall command runs the findcall analyzer.
package main

import (
	"utilware/dep/x/tools/go/analysis/passes/findcall"
	"utilware/dep/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(findcall.Analyzer) }
