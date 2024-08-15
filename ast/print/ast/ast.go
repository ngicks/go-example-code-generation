package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

const src = `package target

import "fmt"

type Foo string

const (
	FooFoo Foo = "foo"
	FooBar Foo = "bar"
	FooBaz Foo = "baz"
)

func Bar(x, y string) string {
	if len(x) == 0 {
		return y + y
	}
	return fmt.Sprintf("%q%q", x, y)
}

type Some[T, U any] struct {
	Foo string
	Bar T
	Baz U
}

func (s Some[T, U]) Method1() {
	// ...nothing...
}

// comment slash slash


/* 

comment slash star

*/
`

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "./target/foo.go", src, parser.AllErrors|parser.ParseComments)
	if err != nil {
		panic(err)
	}
	_ = ast.Print(fset, f)
}
