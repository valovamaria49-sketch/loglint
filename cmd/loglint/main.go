package main

import (
	"github.com/valovamaria49-sketch/loglint/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analysis.Analyzer)
}
