package main

import (
	"github.com/AlexBlackNn/metrics/analyzer/checker"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
)

func main() {
	multichecker.Main(
		checker.CheckAnalyzer, // или errcheckanalyzer.ErrCheckAnalyzer, если анализатор импортируется
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	)
}
