package main

import (
	"fmt"
	"os"
	"text/template"
)

var (
	example = template.Must(template.New("").Parse(
		`Hi {{.Gopher}}.
{{range $idx, $el := .Iter}}    {{if not .}}Hey {{$.Gopher}} this is empty
	{{- if not $.Continue}}{{break}}{{end -}}
{{else}}Iterating at {{$idx}}: {{.Field}} {{end}}
{{end}}
`,
	))
)

func main() {
	decoratingExecute := func(data any) {
		fmt.Println("---")
		err := example.Execute(os.Stdout, data)
		fmt.Println("---")
		fmt.Printf("error: %v\n", err)
		fmt.Println()
	}

	decoratingExecute(map[string]any{
		"Gopher": "you",
		"Iter":   []map[string]string{{"Field": "foo"}, {"Field": "bar"}, {}, {"Field": "baz"}},
	})

	decoratingExecute(map[string]any{
		"Gopher":   "you",
		"Continue": "ok",
		"Iter":     []map[string]string{{"Field": "foo"}, {"Field": "bar"}, {}, {"Field": "baz"}},
	})

	decoratingExecute(map[string]any{
		"Gopher": "you",
		"Iter":   map[string]map[string]string{"0": {"Field": "foo"}, "1": {"Field": "bar"}, "2": {"Field": "baz"}},
	})
}
