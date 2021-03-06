package alphabetical

import (
	"errors"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
)

const targetText = "// Alphabetical order"

var Doc = `sort by alphabetical
// Alphabetical order
above comment check.
`

var Analyzer = &analysis.Analyzer{
	Name:     "alphabetical",
	Doc:      Doc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

var ErrNotAlphabeticalOrder = errors.New("not sort by alphabetical")

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	var (
		commentMap ast.CommentMap
		nowPackage string
	)
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.BlockStmt)(nil),
		(*ast.GenDecl)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch v := n.(type) {
		case *ast.File:
			commentMap = ast.NewCommentMap(pass.Fset, v, v.Comments)
			nowPackage = v.Name.String()
		case *ast.GenDecl:
			checkGenDcl(pass, v)
		case *ast.BlockStmt:
			checkBlock(pass, v, commentMap, nowPackage)
		}
	})

	return nil, nil
}

func checkBlock(pass *analysis.Pass, block *ast.BlockStmt, commentMap ast.CommentMap, nowPackage string) {
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

			fn, args := callName(pass, call, nowPackage)
			if beforeFunc != "" && beforeFunc != fn && beforeName != args {
				checking = false
				continue
			}
			if beforeFunc == fn && beforeName > args {
				pass.Reportf(call.Pos(), ErrNotAlphabeticalOrder.Error())
			}
			if beforeName == args && beforeFunc > fn {
				pass.Reportf(call.Pos(), ErrNotAlphabeticalOrder.Error())
			}
			beforeFunc, beforeName = fn, args
		default:
			checking = false
		}
	}
}

func callName(pass *analysis.Pass, call *ast.CallExpr, nowPackage string) (funcName string, args string) {
	fn := typeutil.StaticCallee(pass.TypesInfo, call)
	if fn == nil {
		return "", ""
	}
	for _, arg := range call.Args {
		switch v := arg.(type) {
		case *ast.BasicLit:
			args += strings.Trim(v.Value, "\"/")
		case *ast.Ident:
			if v.Name == "nil" {
				args += " "
			} else {
				args += v.Name
			}
		case *ast.CallExpr:
			fn, as := callName(pass, v, nowPackage)
			args += fn + as
		}
	}

	if fn.Pkg().Name() == nowPackage {
		return fn.Name(), args
	}
	fsplit := strings.Split(fn.FullName(), "/")
	return fsplit[len(fsplit)-1], args
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

func checkGenDcl(pass *analysis.Pass, gendcl *ast.GenDecl) {
	if gendcl.Doc == nil {
		return
	}
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
				pass.Reportf(spec.Pos(), ErrNotAlphabeticalOrder.Error())
			}
			beforeName = v.Names[0].Name
		case *ast.TypeSpec:
			if beforeName > v.Name.Name {
				pass.Reportf(spec.Pos(), ErrNotAlphabeticalOrder.Error())
			}
			beforeName = v.Name.Name
		}
	}
}
