package main

import (
	"github.com/komem3/alphabetical"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(alphabetical.Analyzer)
}
