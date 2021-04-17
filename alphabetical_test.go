package alphabetical_test

import (
	"testing"

	"github.com/komem3/alphabetical"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, alphabetical.Analyzer, "a")
}
