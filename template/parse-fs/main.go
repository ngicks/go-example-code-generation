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
	baseNameCutExt := func(p string) string {
		p, _ = strings.CutSuffix(path.Base(p), path.Ext(p))
		return p
	}
	for _, tmpl := range tmpls {
		if tmpl.IsDir() {
			continue
		}
		if extTrimmed == nil {
			extTrimmed = template.New(baseNameCutExt(tmpl.Name()))
		}
		bin, err := templates.ReadFile(path.Join("template", tmpl.Name()))
		if err != nil {
			panic(err)
		}
		_ = template.Must(extTrimmed.New(baseNameCutExt(tmpl.Name())).Parse(string(bin)))
	}
}

type param struct {
	Tmp2, Tmp3, Tmp4 sub
}
type sub struct {
	Yay string
	Nay string
}

func main() {
	root = root.Lookup("tmp1.tmpl")
	data := param{
		Tmp2: sub{
			Yay: "yay2",
			Nay: "nay2",
		},
		Tmp3: sub{
			Yay: "yay3",
			Nay: "nay3",
		},
		Tmp4: sub{
			Yay: "yay4",
			Nay: "nay4",
		},
	}
	fmt.Println("---")
	err := root.Execute(os.Stdout, data)
	fmt.Println("---")
	fmt.Printf("err: %v\n", err)
	/*
		---
		tmp2: yay2
		tmp3: yay3
		tmp4: yay4

		---
		err: <nil>
	*/

	_, _ = root.New("additional").Parse(`{{.Tmp2.Yay}} and {{.Tmp3.Nay}}`)

	fmt.Println()
	fmt.Println("---")
	err = root.Execute(os.Stdout, data)
	fmt.Println("---")
	fmt.Printf("err: %v\n", err)
	/*
		---
		tmp2: yay2
		tmp3: yay3
		tmp4: yay4
		yay2 and nay3
		---
		err: <nil>
	*/
}
