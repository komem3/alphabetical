package main

import (
	"github.com/komem3/alphabeticalorder"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(alphabeticalorder.Analyzer)
}
