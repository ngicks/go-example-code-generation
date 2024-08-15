package main

import (
	"errors"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
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

	generatedDir := filepath.Join("ast", "rewrite", "dstutil", "generated")
	err = os.Mkdir(generatedDir, fs.ModePerm)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		panic(err)
	}

	pkg := pkgs[0]

	for _, f := range pkg.Syntax {
		df, err := decorator.DecorateFile(pkg.Fset, f)
		if err != nil {
			panic(err)
		}
		dstutil.Apply(
			df,
			func(c *dstutil.Cursor) bool {
				n := c.Node()
				switch x := n.(type) {
				default:
					return true
				case *dst.FuncDecl:
				case *dst.GenDecl:
					if x.Tok != token.TYPE {
						break
					}

					if len(x.Specs) == 1 {
						name, ok := isStringBasedType(x.Specs[0])
						if !ok {
							break
						}
						param, ok := parseDirective(x.Decs)
						if !ok {
							break
						}
						param.Name = name
						addOrReplaceEnum(c, param)
					}
				}
				return false
			},
			nil,
		)

		out, err := os.Create(filepath.Join(generatedDir, filepath.Base(pkg.Fset.Position(f.FileStart).Filename)))
		if err != nil {
			panic(err)
		}

		restorer := decorator.NewRestorer()
		af, err := restorer.RestoreFile(df)
		if err != nil {
			panic(err)
		}
		printer.Fprint(out, restorer.Fset, af)
	}
}

func isStringBasedType(spec dst.Spec) (string, bool) {
	typ, ok := spec.(*dst.TypeSpec)
	if !ok {
		return "", false
	}
	id, ok := typ.Type.(*dst.Ident)
	if !ok {
		return "", false
	}
	return typ.Name.Name, id.Name == "string"
}

type EnumParam struct {
	Name     string
	Variants []string
}

func parseDirective(decorations dst.GenDeclDecorations) (EnumParam, bool) {
	for i := len(decorations.Start) - 1; i >= 0; i-- {
		line := decorations.Start[i]
		if len(strings.TrimSpace(line)) == 0 {
			// start of comments groups that is not associated to gen decl.
			break
		}
		c := stripMarker(line)
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

func addOrReplaceEnum(c *dstutil.Cursor, param EnumParam) {
	found := false
	dstutil.Apply(
		c.Parent(),
		func(c *dstutil.Cursor) bool {
			node := c.Node()
			switch x := node.(type) {
			default:
				return true
			case *dst.FuncDecl:
			case *dst.GenDecl:
				if x.Tok != token.CONST {
					break
				}
				if !isGeneratedFor(x.Decs, param.Name) {
					break
				}
				found = true
				c.Replace(astVariants(param, x.Decs))
			}
			return false
		},
		nil,
	)
	if !found {
		c.InsertAfter(astVariants(param, dst.GenDeclDecorations{}))
	}
}

func isGeneratedFor(decorations dst.GenDeclDecorations, fotTy string) bool {
	for i := len(decorations.Start) - 1; i >= 0; i-- {
		line := decorations.Start[i]
		if len(strings.TrimSpace(line)) == 0 {
			break
		}
		c := stripMarker(line)
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

func astVariants(param EnumParam, targetDecoration dst.GenDeclDecorations) *dst.GenDecl {
	if len(targetDecoration.Start) > 0 {
		var i int
		for i = len(targetDecoration.Start) - 1; i >= 0; i-- {
			if targetDecoration.Start[i] == "\n" {
				break
			}
		}
		if i < 0 {
			i = len(targetDecoration.Start) - 1
		}
		targetDecoration.Start = append(slices.Clone(targetDecoration.Start[:i]), "\n", "//enum:generated_for="+param.Name)
	} else {
		targetDecoration.Start = []string{"//enum:generated_for=" + param.Name}
	}
	return &dst.GenDecl{
		Decs:   targetDecoration,
		Tok:    token.CONST,
		Lparen: true,
		Specs:  mapParamToSpec(param),
		Rparen: true,
	}
}

func mapParamToSpec(param EnumParam) []dst.Spec {
	specs := make([]dst.Spec, len(param.Variants))
	for i, variant := range param.Variants {
		specs[i] = &dst.ValueSpec{
			Names:  []*dst.Ident{{Name: param.Name + capitalize(variant)}},
			Type:   &dst.Ident{Name: param.Name},
			Values: []dst.Expr{&dst.BasicLit{Kind: token.STRING, Value: strconv.Quote(variant)}},
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
