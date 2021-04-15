package alphabetorder

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const targetText = "// Alphabetical order"

var Doc = `sort variable by alphabet order
// Alphabetical order
above comment check.
`

var Analyzer = &analysis.Analyzer{
	Name:     "alphabetorder",
	Doc:      Doc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		gendcl := n.(*ast.GenDecl)

		var hit bool
		for _, c := range gendcl.Doc.List {
			if c.Text == targetText {
				hit = true
				break
			}
		}
		if !hit {
			return
		}

		var beforeName string
		for _, spec := range gendcl.Specs {
			switch v := spec.(type) {
			case *ast.ValueSpec:
				if beforeName > v.Names[0].Name {
					pass.Reportf(spec.Pos(), "not sort by alpahbet order")
					return
				}
				beforeName = v.Names[0].Name
			case *ast.TypeSpec:
				if beforeName > v.Name.Name {
					pass.Reportf(spec.Pos(), "not sort by alpahbet order")
					return
				}
				beforeName = v.Name.Name
			}
		}
	})

	return nil, nil
}
