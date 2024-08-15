package main

import (
	"context"
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/tools/go/packages"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedExportFile |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes |
			packages.NeedModule |
			packages.NeedEmbedFiles |
			packages.NeedEmbedPatterns,
		Context: ctx,
		Logf: func(format string, args ...interface{}) {
			fmt.Printf("log: "+format, args...)
			fmt.Println()
		},
	}
	pkgs, err := packages.Load(cfg, "io", "./ast/parse-by-packages/target")
	if err != nil {
		panic(err)
	}
	packages.Visit(
		pkgs,
		func(p *packages.Package) bool {
			for _, err := range p.Errors {
				fmt.Printf("pkg %s: %#v\n", p.PkgPath, err)
			}
			if len(p.Errors) > 0 {
				return true
			}

			fmt.Printf("package path: %s\n", p.PkgPath)

			return true
		},
		nil,
	)

	ioPkg := pkgs[0]
	targetPkg := pkgs[1]

	for _, f := range targetPkg.Syntax {
		ast.Print(targetPkg.Fset, f)
		fmt.Printf("\n\n")
	}

	foo := targetPkg.Types.Scope().Lookup("Foo")
	fmt.Printf("foo: %#v\n", foo)

	r := ioPkg.Types.Scope().Lookup("Reader")

	pfoo := types.NewPointer(foo.Type())

	if types.AssignableTo(pfoo, r.Type()) {
		fmt.Printf("%s satisfies %s\n", pfoo, r)
	}

	mset := types.NewMethodSet(pfoo)
	for i := 0; i < mset.Len(); i++ {
		meth := mset.At(i).Obj()
		sig := meth.Type().Underlying().(*types.Signature)
		fmt.Printf(
			"%d: func (receiver=%s name=*%s)(func name=%s)(params=%s) (results=%s)\n",
			i, sig.Recv().Name(), foo.Name(), meth.Name(), sig.Params(), sig.Results(),
		)
	}
}
