package main

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/dave/jennifer/jen"
)

func main() {
	decoratePrint := func(v any) {
		fmt.Println("---")
		fmt.Printf("%#v\n", v)
		fmt.Println("---")
		fmt.Println()
	}

	decoratePrint(jen.Var().Id("yay").Op("=").Lit("yay yay"))
	/*
		---
		var yay = "yay yay"
		---
	*/

	var f *jen.File

	f = jen.NewFile("baz")
	f.NoFormat = true
	f.Qual("foo/bar/baz", "Wow")
	decoratePrint(f)
	/*
		---
		package baz

		import baz "foo/bar/baz"


		baz.Wow
		---
	*/

	f = jen.NewFilePath("foo/bar/baz")
	f.NoFormat = true
	f.Qual("foo/bar/baz", "Wow")
	decoratePrint(f)
	/*
		---
		package baz


		Wow
		---
	*/

	f = jen.NewFilePathName("foo/bar/baz", "hoge")
	f.NoFormat = true
	f.Qual("foo/bar/baz", "Wow")
	decoratePrint(f)
	/*
	   ---
	   package hoge


	   Wow
	   ---
	*/

	f = jen.NewFile("foo")
	f.PackageComment("// Code generated by me. DO NOT EDIT.")
	decoratePrint(f)
	/*
		---
		// Code generated by me. DO NOT EDIT.
		package foo

		---
	*/

	decoratePrint(jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())))
	/*
	   ---
	   if err != nil {
	           return err
	   }
	   ---
	*/

	decoratePrint(jen.Type().Id("foo").Struct(
		jen.Id("A").String().Tag(map[string]string{"json": "a"}),
		jen.Id("B").Int().Tag(map[string]string{"json": "b", "bar": "baz"}),
	))
	/*
		---
		type foo struct {
		        A string `json:"a"`
		        B int    `bar:"baz" json:"b"`
		}
		---
	*/

	decoratePrint(jen.Var().Id("bar").Op("=").Index(jen.Op("...")).String().Values(jen.Lit("foo"), jen.Lit("bar"), jen.Lit("baz")))
	/*
	   ---
	   var bar = [...]string{"foo", "bar", "baz"}
	   ---
	*/

	decoratePrint(
		jen.Var().Id("bar").Op("=").Index(jen.Op("...")).String().
			Custom(jen.Options{
				Open:      "{",
				Close:     "}",
				Separator: ",",
				Multi:     true,
			},
				jen.Lit("foo"), jen.Lit("bar"), jen.Lit("baz"),
			),
	)
	/*
		---
		var bar = [...]string{
		        "foo",
		        "bar",
		        "baz",
		}
		---
	*/

	fields := []struct {
		name string
		def  *jen.Statement
	}{
		{"foo", jen.String()},
		{"bar", jen.Int()},
		{"baz", jen.Op("*").Qual("bytes", "Buffer")},
	}
	decoratePrint(jen.Type().Id("foo").StructFunc(func(g *jen.Group) {
		for _, f := range fields {
			g.Id(f.name).Add(f.def)
		}
	}))
	/*
		---
		type foo struct {
		        foo string
		        bar int
		        baz *bytes.Buffer
		}
		---
	*/

	decoratePrint(jen.Type().Id("bar").Op("struct").Op("{").
		Do(func(s *jen.Statement) {
			for _, f := range fields {
				s.Id(f.name).Add(f.def).Line()
			}
		}).
		Op("}"),
	)
	/*
		---
		type bar struct {
		        foo string
		        bar int
		        baz *bytes.Buffer
		}
		---
	*/

	randomIdentName := func(leng int) string {
		var buf strings.Builder
		buf.Grow(1 + 2*leng)
		_ = buf.WriteByte('_')
		for range leng {
			_, _ = buf.WriteString(fmt.Sprintf("%x", [1]byte{rand.N[byte](255)}))
		}
		return buf.String()
	}
	bytesQual := func(ident string) *jen.Statement {
		return jen.Qual("bytes", ident)
	}
	imageQual := func(ident string) *jen.Statement {
		return jen.Qual("image", ident)
	}
	decoratePrint(jen.Var().DefsFunc(func(g *jen.Group) {
		g.Id(randomIdentName(4)).Op("=").Op("*").Add(bytesQual("Buffer"))
		g.Id(randomIdentName(4)).Op("=").Add(imageQual("Image"))
	}))
	/*
		---
		var (
		        _76ef57e5 = *bytes.Buffer
		        _5d5cdb42 = image.Image
		)
		---
	*/

	decoratePrint(jen.Type().Id("Yay").String().Line().Id(`func foo() bool { return true }`).Line().Type().Id("Nay").Bool())
	/*
		type Yay string

		func foo() bool { return true }

		type Nay bool
	*/
}
