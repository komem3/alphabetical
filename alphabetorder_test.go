package alphabetorder_test

import (
	"testing"

	"github.com/komem3/alphabetorder"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, alphabetorder.Analyzer, "a")
}
