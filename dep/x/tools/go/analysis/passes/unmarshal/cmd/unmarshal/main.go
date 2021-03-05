// The unmarshal command runs the unmarshal analyzer.
package main

import (
	"utilware/dep/x/tools/go/analysis/passes/unmarshal"
	"utilware/dep/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(unmarshal.Analyzer) }
