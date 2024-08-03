package main

import (
	"embed"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"
)

//go:embed template
var templates embed.FS

var (
	root       = template.Must(template.ParseFS(templates, "template/*"))
	extTrimmed *template.Template
)

func init() {
	tmpls, err := templates.ReadDir("template")
	if err != nil {
		panic(err)
	}
	cutExt := func(p string) string {
		p, _ = strings.CutSuffix(path.Base(p), path.Ext(p))
		return p
	}
	for _, tmpl := range tmpls {
		if tmpl.IsDir() {
			continue
		}
		if extTrimmed == nil {
			extTrimmed = template.New(cutExt(tmpl.Name()))
		}
		bin, err := templates.ReadFile(path.Join("template", tmpl.Name()))
		if err != nil {
			panic(err)
		}
		_ = template.Must(extTrimmed.New(cutExt(tmpl.Name())).Parse(string(bin)))
	}
}

type param struct {
	Sub1, Sub2, Sub3 sub
}
type sub struct {
	Yay string
	Nay string
}

func main() {
	root = root.Lookup("root.tmpl")
	data := param{
		Sub1: sub{
			Yay: "yay1",
			Nay: "nay1",
		},
		Sub2: sub{
			Yay: "yay2",
			Nay: "nay2",
		},
		Sub3: sub{
			Yay: "yay3",
			Nay: "nay3",
		},
	}
	fmt.Println("---")
	err := root.Execute(os.Stdout, data)
	fmt.Println("---")
	fmt.Printf("err: %v\n", err)
	/*
		---
		sub1: yay1
		sub2: yay2
		sub3: yay3

		---
		err: <nil>
	*/

	_, _ = root.New("additional").Parse(`{{.Sub1.Yay}} and {{.Sub2.Nay}}`)

	fmt.Println()
	fmt.Println("---")
	err = root.Execute(os.Stdout, data)
	fmt.Println("---")
	fmt.Printf("err: %v\n", err)
	/*
		---
		sub1: yay1
		sub2: yay2
		sub3: yay3
		yay1 and nay2
		---
		err: <nil>
	*/
}
