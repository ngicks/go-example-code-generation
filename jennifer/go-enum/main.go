package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/dave/jennifer/jen"
)

type EnumParam struct {
	PackageName string
	Name        string
	Variants    []string
	Excepts     []EnumExceptParam
}

type EnumExceptParam struct {
	Name             string
	ExceptName       string
	ExcludedValiants []string
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

func replaceInvalidChar(s string) string {
	// As per Go programming specification.
	// identifier = letter { letter | unicode_digit }.
	// https://go.dev/ref/spec#Identifiers
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r) {
			return r
		}
		return '_'
	}, s)
}

func fillName(p EnumExceptParam, name string) EnumExceptParam {
	p.Name = name
	return p
}

func main() {
	pkgPath := filepath.Join("jennifer", "go-enum", "example")
	err := os.MkdirAll(pkgPath, fs.ModePerm)
	if err != nil {
		panic(err)
	}

	param := EnumParam{
		PackageName: "example",
		Name:        "Enum",
		Variants:    []string{"foo", "b\"ar", "baz"},
		Excepts: []EnumExceptParam{
			{
				ExceptName:       "foo",
				ExcludedValiants: []string{"foo"},
			},
			{
				ExceptName:       "Muh",
				ExcludedValiants: []string{"foo", "b\"ar"},
			},
		},
	}

	f := jen.NewFile(param.PackageName)

	f.PackageComment("// Code generated by me. DO NOT EDIT.")

	out, err := os.Create(filepath.Join(pkgPath, "enum.go"))
	if err != nil {
		panic(err)
	}

	f.Type().Id(param.Name).String() // type Enum string

	// const (
	f.Const().DefsFunc(func(g *jen.Group) {
		for _, variant := range param.Variants {
			g.
				Id(param.Name + replaceInvalidChar(capitalize(variant))). // EnumFoo
				Id(param.Name).                                           // Enum
				Op("=").                                                  // =
				Lit(variant)                                              // "foo"\n
		}
	}) // )

	// var _EnumAll = [...]Enum
	f.Var().Id("_" + param.Name + "All").Op("=").Index(jen.Op("...")).Id(param.Name).
		ValuesFunc(func(g *jen.Group) { // {
			for _, variant := range param.Variants {
				g.Line().Id(param.Name + replaceInvalidChar(capitalize(variant))) // EnumFoo,
			}
			g.Line() // \n
		}) // }

	// func IsEnum(v Enum) bool
	f.Func().Id("Is" + param.Name).Params(jen.Id("v").Id(param.Name)).Bool().Block( // {
		jen.Return( // return
			jen.Qual("slices", "Contains").Call( // slices.Contains
				jen.Id("_"+param.Name+"All").Index(jen.Op(":")), // _EnumAll[:],
				jen.Id("v"), // v,
			),
		),
	) // }

	f.Line()

	for _, except := range param.Excepts {
		except = fillName(except, param.Name)
		// func IsEnumExceptFoo(v Enum) bool
		f.Func().Id("Is" + except.Name + "Except" + replaceInvalidChar(capitalize(except.ExceptName))).Params(jen.Id("v").Id(except.Name)).Bool().Block( // {
			jen.Return( // return
				jen.Op("!").Qual("slices", "Contains").Params( // !slice.Contains(
					jen.Line().Index().Id(except.Name).ValuesFunc(func(g *jen.Group) { //[]Enum{
						for _, e := range except.ExcludedValiants {
							g.Line().Id(param.Name + replaceInvalidChar(capitalize(e))) // EnumFoo,
						}
						g.Line()
					}), // },
					jen.Line().Id("v"), // v,
					jen.Line(),
				), // )
			),
		) // }
		f.Line()
	}

	err = f.Render(out)
	if err != nil {
		panic(err)
	}
}
