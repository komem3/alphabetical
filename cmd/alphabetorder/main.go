package main

import (
	"github.com/komem3/alphabetorder"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(alphabetorder.Analyzer)
}
