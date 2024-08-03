package main

import (
	"errors"
	"fmt"
	"os"
	"text/template"
)

var (
	example = template.Must(template.New("").Parse(
		`An example template.
Hello {{.Gopher}}.
Yay Yay.
`,
	))
	chained = template.Must(template.New("").Parse(`**chained**{{.Chain.Gopher}}
`))
	accessingUnexported = template.Must(template.New("").Parse(
		`accessing unexported field: {{.unexportedField}}
`,
	))
)

type sample struct {
	Gopher string
}

type embedded struct {
	sample
}

type sampleMethod1 struct {
}

func (s sampleMethod1) Gopher() string {
	return "method"
}

type sampleMethod2 struct {
	err error
}

func (s sampleMethod2) Gopher() (string, error) {
	return "method2", s.err
}

type chainedData struct {
	v   any
	err error
}

func (c chainedData) Chain() (any, error) {
	return c.v, c.err
}

type unexported struct {
	unexportedField string
}

func main() {
	decoratingExecute := func(t *template.Template, data any) {
		fmt.Println("---")
		err := t.Execute(os.Stdout, data)
		fmt.Println("---")
		fmt.Printf("error: %v\n", err)
		fmt.Println()
	}

	executeTemplate := func(data any) { decoratingExecute(example, data) }
	executeTemplate(sample{Gopher: "me"})
	executeTemplate(sample{Gopher: "you"})
	executeTemplate(embedded{sample{Gopher: "embedded"}})
	executeTemplate(map[string]string{"Gopher": "from map[string]string"})
	executeTemplate(map[string]any{"Gopher": "from map[string]any"})
	executeTemplate(sampleMethod1{})
	executeTemplate(sampleMethod2{})
	executeTemplate(sampleMethod2{err: errors.New("sample")})

	executeChained := func(data any) { decoratingExecute(chained, data) }
	executeChained(chainedData{v: sampleMethod2{}})
	executeChained(chainedData{v: map[string]string{"Gopher": "map"}})

	executeAccessingUnexported := func(data any) { decoratingExecute(accessingUnexported, data) }
	executeAccessingUnexported(map[string]string{"unexportedField": "unexported field"})
	executeAccessingUnexported(unexported{})
}
