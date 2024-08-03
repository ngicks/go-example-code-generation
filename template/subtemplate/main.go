package main

import (
	"fmt"
	"os"
	"text/template"
)

var (
	root = template.Must(template.New("root").Parse(
		`sub1: {{template "sub1" .Sub1}}
sub2: {{template "sub2" .Sub2}}
sub3: {{template "sub3" .Sub3}}
{{block "additional" .}}{{end}}
`))
	sub1 = template.Must(root.New("sub1").Parse(`{{.Yay}}`))
	_    = template.Must(root.New("sub2").Parse(`{{.Yay}}`))
	sub3 = template.Must(root.New("sub3").Parse(`{{.Yay}}`))
)

type param struct {
	Sub1, Sub2, Sub3 sub
}
type sub struct {
	Yay string
	Nay string
}

func main() {
	decoratingExecute := func(data any) {
		fmt.Println("---")
		err := root.Execute(os.Stdout, data)
		fmt.Println("---")
		fmt.Printf("error: %v\n", err)
		fmt.Println()
	}

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

	decoratingExecute(data)
	/*
		---
		sub1: yay1
		sub2: yay2
		sub3: yay3

		---
		error: <nil>
	*/
	_, _ = sub1.Parse(`{{.Nay}}`)
	decoratingExecute(data)
	/*
		---
		sub1: nay1
		sub2: yay2
		sub3: yay3

		---
		error: <nil>
	*/

	_, _ = sub3.New("additional").Parse(`{{.Sub1.Yay}} and {{.Sub2.Nay}}`)
	decoratingExecute(data)
	/*
		---
		sub1: nay1
		sub2: yay2
		sub3: yay3
		yay1 and nay2
		---
		error: <nil>
	*/
}
