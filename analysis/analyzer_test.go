package analysis

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzerOnTestdata(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "logs", "zaplogs")
}
