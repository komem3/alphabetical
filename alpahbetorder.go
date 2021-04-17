package alphabetorder

import (
	"errors"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
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

var ErrNotAlphabetOrder = errors.New("not sort by alphabet order")

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	var commentMap ast.CommentMap
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.BlockStmt)(nil),
		(*ast.GenDecl)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		var (
			err error
			pos token.Pos
		)
		switch v := n.(type) {
		case *ast.File:
			commentMap = ast.NewCommentMap(pass.Fset, v, v.Comments)
		case *ast.GenDecl:
			pos, err = checkGenDcl(v)
		case *ast.BlockStmt:
			pos, err = checkBlock(pass, v, commentMap)
		}

		if err != nil {
			pass.Reportf(pos, err.Error())
		}
	})

	return nil, nil
}

func checkBlock(pass *analysis.Pass, block *ast.BlockStmt, commentMap ast.CommentMap) (token.Pos, error) {
	var (
		checking   bool
		beforeName string
		beforeFunc string
	)
	if checkComment(block, commentMap) {
		checking, beforeName, beforeFunc = true, "", ""
	}
	for _, stmt := range block.List {
		if checkComment(stmt, commentMap) {
			checking, beforeName, beforeFunc = true, "", ""
		}
		if !checking {
			continue
		}

		switch v := stmt.(type) {
		case *ast.ExprStmt:
			call, ok := v.X.(*ast.CallExpr)
			if !ok {
				checking = false
				continue
			}

			fn, args := callName(pass, call)
			if fn == "" || (beforeFunc != "" && beforeFunc != fn) {
				checking = false
				continue
			}
			if beforeName > args {
				return call.Pos(), ErrNotAlphabetOrder
			}
			beforeFunc, beforeName = fn, args
		}
	}
	return 0, nil
}

func callName(pass *analysis.Pass, call *ast.CallExpr) (funcName string, args string) {
	fn := typeutil.StaticCallee(pass.TypesInfo, call)
	if fn == nil {
		return "", ""
	}
	for _, arg := range call.Args {
		switch v := arg.(type) {
		case *ast.BasicLit:
			args += strings.Trim(v.Value, "\"/")
		case *ast.Ident:
			if v.Name != "nil" {
				args += v.Name
			}
		case *ast.CallExpr:
			fn, as := callName(pass, v)
			args += fn + as
		}
	}
	return fn.FullName(), args
}

func checkComment(n ast.Node, commentMap ast.CommentMap) bool {
	if comments, ok := commentMap[n]; ok {
		for _, commentGroup := range comments {
			for _, comment := range commentGroup.List {
				if comment.Text == targetText {
					return true
				}
			}
		}
	}
	return false
}

func checkGenDcl(gendcl *ast.GenDecl) (token.Pos, error) {
	if gendcl.Doc == nil {
		return 0, nil
	}
	var hit bool
	for _, c := range gendcl.Doc.List {
		if c.Text == targetText {
			hit = true
			break
		}
	}
	if !hit {
		return 0, nil
	}

	var beforeName string
	for _, spec := range gendcl.Specs {
		switch v := spec.(type) {
		case *ast.ValueSpec:
			if beforeName > v.Names[0].Name {
				return spec.Pos(), ErrNotAlphabetOrder
			}
			beforeName = v.Names[0].Name
		case *ast.TypeSpec:
			if beforeName > v.Name.Name {
				return spec.Pos(), ErrNotAlphabetOrder
			}
			beforeName = v.Name.Name
		}
	}
	return 0, nil
}
