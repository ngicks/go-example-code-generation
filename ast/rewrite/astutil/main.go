package main

import (
	"errors"
	"go/ast"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

func main() {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedImports |
			packages.NeedTypes |
			packages.NeedSyntax,
	}
	pkgs, err := packages.Load(cfg, "./ast/rewrite/target")
	if err != nil {
		panic(err)
	}

	generatedDir := filepath.Join("ast", "rewrite", "astutil", "generated")
	err = os.Mkdir(generatedDir, fs.ModePerm)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		panic(err)
	}

	pkg := pkgs[0]

	for _, f := range pkg.Syntax {
		cm := ast.NewCommentMap(pkg.Fset, f, f.Comments)
		filename := filepath.Base(pkg.Fset.Position(f.FileStart).Filename)
		astutil.Apply(
			f,
			func(c *astutil.Cursor) bool {
				n := c.Node()
				switch x := n.(type) {
				default:
					return true
				case *ast.FuncDecl:
				case *ast.GenDecl:
					if x.Tok != token.TYPE {
						break
					}

					if len(x.Specs) == 1 {
						name, ok := isStringBasedType(x.Specs[0])
						if !ok {
							break
						}
						param, ok := parseDirective(x.Doc)
						if !ok {
							break
						}
						param.Name = name
						addOrReplaceEnum(c, param, cm)
					}
				}
				return false
			},
			nil,
		)

		f.Comments = cm.Comments()

		out, err := os.Create(filepath.Join(generatedDir, filename))
		if err != nil {
			panic(err)
		}

		err = printer.Fprint(out, pkg.Fset, f)
		if err != nil {
			panic(err)
		}
	}
}

func isStringBasedType(spec ast.Spec) (string, bool) {
	typ, ok := spec.(*ast.TypeSpec)
	if !ok {
		return "", false
	}
	id, ok := typ.Type.(*ast.Ident)
	if !ok {
		return "", false
	}
	return typ.Name.Name, id.Name == "string"
}

type EnumParam struct {
	Name     string
	Variants []string
}

func parseDirective(cg *ast.CommentGroup) (EnumParam, bool) {
	for _, comment := range cg.List {
		c := strings.TrimLeftFunc(stripMarker(comment.Text), unicode.IsSpace)
		c, isDirection := strings.CutPrefix(c, "enum:variants=")
		if !isDirection {
			continue
		}
		return EnumParam{Variants: strings.Split(c, ",")}, true
	}
	return EnumParam{}, false
}

func stripMarker(text string) string {
	if len(text) < 2 {
		return text
	}
	switch text[1] {
	case '/':
		return text[2:]
	case '*':
		return text[2 : len(text)-2]
	}
	return text
}

func addOrReplaceEnum(c *astutil.Cursor, param EnumParam, cm ast.CommentMap) {
	found := false
	astutil.Apply(
		c.Parent(),
		func(c *astutil.Cursor) bool {
			node := c.Node()
			switch x := node.(type) {
			default:
				return true
			case *ast.FuncDecl:
			case *ast.GenDecl:
				if x.Tok != token.CONST {
					break
				}
				if !isGeneratedFor(x.Doc, param.Name) {
					break
				}
				found = true
				newDecl := astVariants(param, x.TokPos)
				delete(cm, x)
				cm[newDecl] = append(cm[newDecl], newDecl.Doc)
				c.Replace(newDecl)
			}
			return false
		},
		nil,
	)
	if !found {
		newDecl := astVariants(param, c.Node().(*ast.GenDecl).Specs[0].Pos()+30)
		cm[newDecl] = append(cm[newDecl], newDecl.Doc)
		c.InsertAfter(newDecl)
	}
}

func isGeneratedFor(cg *ast.CommentGroup, fotTy string) bool {
	for _, comment := range cg.List {
		c := strings.TrimLeftFunc(stripMarker(comment.Text), unicode.IsSpace)
		s, ok := strings.CutPrefix(c, "enum:generated_for=")
		if !ok {
			return false
		}
		if s == fotTy {
			return true
		}
	}
	return false
}

func astVariants(param EnumParam, pos token.Pos) *ast.GenDecl {
	return &ast.GenDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Slash: pos,
					Text:  "//enum:generated_for=" + param.Name,
				},
			},
		},
		TokPos: token.Pos(int(pos) + len("//enum:generated_for="+param.Name) + 1),
		Tok:    token.CONST,
		Lparen: 1,
		Specs:  mapParamToSpec(param),
		Rparen: 2,
	}
}

func mapParamToSpec(param EnumParam) []ast.Spec {
	specs := make([]ast.Spec, len(param.Variants))
	for i, variant := range param.Variants {
		specs[i] = &ast.ValueSpec{
			Names:  []*ast.Ident{{Name: param.Name + capitalize(variant)}},
			Type:   &ast.Ident{Name: param.Name},
			Values: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(variant)}},
		}
	}
	return specs
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	if len(s) == 1 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
