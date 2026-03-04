package exit

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestNoExit(t *testing.T) {
	analysistest.Run(t, analysistest.TestData()+"/no_exit", Analyzer)
}

func TestNotMain(t *testing.T) {
	analysistest.Run(t, analysistest.TestData()+"/not_main", Analyzer)
}

func TestHasExit(t *testing.T) {
	analysistest.Run(t, analysistest.TestData()+"/has_exit", Analyzer)
}
