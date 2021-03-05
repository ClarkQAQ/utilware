// The shadow command runs the shadow analyzer.
package main

import (
	"utilware/dep/x/tools/go/analysis/passes/shadow"
	"utilware/dep/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(shadow.Analyzer) }
