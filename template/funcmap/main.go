package main

import (
	"fmt"
	"os"
	"text/template"
)

var (
	example = template.Must(
		template.
			New("").
			Funcs(template.FuncMap{"customFunc": func() string { return "" }}).
			Parse(
				`{{customFunc .}}
`,
			),
	)
)

func main() {
	decoratingExecute := func(funcs template.FuncMap, data any) {
		fmt.Println("---")
		err := example.Funcs(funcs).Execute(os.Stdout, data)
		fmt.Println("---")
		fmt.Printf("error: %v\n", err)
		fmt.Println()
	}

	decoratingExecute(nil, "foo")
	/*
		---
		---
		error: template: :1:2: executing "" at <customFunc>: wrong number of args for customFunc: want 0 got 1
	*/
	decoratingExecute(
		template.FuncMap{"customFunc": func(v any) string { return fmt.Sprintf("%s", v) }},
		"foo",
	)
	/*
		---
		foo
		---
		error: <nil>
	*/
	decoratingExecute(
		template.FuncMap{"customFunc": func(v ...any) string {
			fmt.Printf("customFunc: %#v\n", v)
			return "ah"
		}},
		"bar",
	)
	/*
		---
		customFunc: []interface {}{"bar"}
		ah
		---
		error: <nil>
	*/
	decoratingExecute(
		template.FuncMap{"customFunc": func(v string) string {
			fmt.Printf("customFunc: %#v\n", v)
			return "ah"
		}},
		"baz",
	)
	/*
		---
		customFunc: "baz"
		ah
		---
		error: <nil>
	*/
	decoratingExecute(
		template.FuncMap{"customFunc": func(v int) string {
			return "ah"
		}},
		"qux",
	)
	/*
		---
		---
		error: template: :1:13: executing "" at <.>: wrong type for value; expected int; got string
	*/
	type sample struct {
		Foo string
		Bar int
	}
	decoratingExecute(
		template.FuncMap{"customFunc": func(v any) int {
			fmt.Printf("customFunc: %#v\n", v)
			return v.(sample).Bar
		}},
		sample{Foo: "foo", Bar: 123},
	)
	/*
	   ---
	   customFunc: main.sample{Foo:"foo", Bar:123}
	   123
	   ---
	   error: <nil>
	*/
	decoratingExecute(
		template.FuncMap{"customFunc": func(v sample) string {
			return v.Foo
		}},
		sample{Foo: "foo", Bar: 123},
	)
	/*
		---
		foo
		---
		error: <nil>
	*/
}
