package alphabeticalorder_test

import (
	"testing"

	"github.com/komem3/alphabeticalorder"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, alphabeticalorder.Analyzer, "a")
}
