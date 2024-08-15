package main

import (
	"fmt"
	"os"
	"text/template"
)

var (
	tmp1 = template.Must(template.New("tmp1").Parse(
		`tmp2: {{template "tmp2" .Tmp2}}
tmp3: {{template "tmp3" .Tmp3}}
tmp4: {{template "tmp4" .Tmp4}}
{{block "additional" .}}{{end}}
`))
	tmp2 = template.Must(tmp1.New("tmp2").Parse(`{{.Yay}}`))
	_    = template.Must(tmp1.New("tmp3").Parse(`{{.Yay}}`))
	tmp4 = template.Must(tmp1.New("tmp4").Parse(`{{.Yay}}`))
)

type param struct {
	Tmp2, Tmp3, Tmp4 sub
}
type sub struct {
	Yay string
	Nay string
}

func main() {
	decoratingExecute := func(data any) {
		fmt.Println("---")
		err := tmp1.Execute(os.Stdout, data)
		fmt.Println("---")
		fmt.Printf("error: %v\n", err)
		fmt.Println()
	}

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

	decoratingExecute(data)
	/*
		---
		tmp2: yay2
		tmp3: yay3
		tmp4: yay4

		---
		error: <nil>
	*/
	_, _ = tmp2.Parse(`{{.Nay}}`)
	decoratingExecute(data)
	/*
		---
		tmp2: nay2
		tmp3: yay3
		tmp4: yay4

		---
		error: <nil>
	*/

	_, _ = tmp4.New("additional").Parse(`{{.Tmp2.Yay}} and {{.Tmp3.Nay}}`)
	decoratingExecute(data)
	/*
		---
		tmp2: nay2
		tmp3: yay3
		tmp4: yay4
		yay2 and nay3
		---
		error: <nil>
	*/
}
