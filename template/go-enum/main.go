package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"unicode"
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

var funcs = template.FuncMap{
	"capitalize": func(s string) string {
		if len(s) == 0 {
			return s
		}
		if len(s) == 1 {
			return strings.ToUpper(s)
		}
		return strings.ToUpper(s[:1]) + s[1:]
	},
	"replaceInvalidChar": func(s string) string {
		// As per Go programming specification.
		// identifier = letter { letter | unicode_digit }.
		// https://go.dev/ref/spec#Identifiers
		return strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r) {
				return r
			}
			return '_'
		}, s)
	},
	"quote": func(s string) string {
		return strconv.Quote(s)
	},
	"fillName": func(p EnumExceptParam, name string) EnumExceptParam {
		p.Name = name
		return p
	},
}

var (
	pkg = template.Must(template.New("package").Funcs(funcs).Parse(
		`// Code generated by me. DO NOT EDIT.
package {{.PackageName}}

import (
	"slices"
)

type {{.Name}} string

const (
{{range .Variants}}	{{$.Name}}{{replaceInvalidChar (capitalize .)}} {{$.Name}} = {{quote .}}
{{end -}}
)

var _{{.Name}}All = [...]{{.Name}}{{"{"}}{{range .Variants}}
	{{$.Name}}{{replaceInvalidChar (capitalize .)}},{{end}}
}

func Is{{.Name}}(v {{.Name}}) bool {
	return slices.Contains(_{{.Name}}All[:], v)
}

{{range .Excepts}}
{{template "except" (fillName . $.Name)}}{{end}}`))
	_ = template.Must(pkg.New("except").Parse(
		`func Is{{.Name}}Except{{replaceInvalidChar (capitalize .ExceptName)}}(v {{.Name}}) bool {
	return !slices.Contains(
		[]{{.Name}}{{"{"}}{{range .ExcludedValiants}}
			{{$.Name}}{{replaceInvalidChar (capitalize .)}},{{end}}
		},
		v,
	)
}
`))
)

func main() {
	pkgPath := filepath.Join("template", "go-enum", "example")
	err := os.MkdirAll(pkgPath, fs.ModePerm)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(filepath.Join(pkgPath, "enum.go"))
	if err != nil {
		panic(err)
	}

	err = pkg.Execute(
		f,
		EnumParam{
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
		},
	)
	if err != nil {
		panic(err)
	}
}
