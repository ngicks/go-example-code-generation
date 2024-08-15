package main

import (
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"os"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
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

// free floating comment

// doc comment
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
	dec := decorator.NewDecorator(fset)
	df, err := dec.DecorateFile(f)
	if err != nil {
		panic(err)
	}
	_ = dst.Print(df)

	dn := dec.Dst.Nodes[f.Decls[0]]

	restorer := decorator.NewRestorer()
	af, err := restorer.RestoreFile(df)
	if err != nil {
		panic(err)
	}
	noop(af)
	an := restorer.Ast.Nodes[dn]
	fmt.Println()
	err = printer.Fprint(os.Stdout, restorer.Fset, an)
	if err != nil {
		panic(err)
	}
	fmt.Println()
}

func noop[T any](_ T) {}
