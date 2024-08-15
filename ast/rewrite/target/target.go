package target

//enum:variants=foo,bar,baz
type Enum string

//enum:variants=foo,bar,baz
type Enum2 string

//enum:generated_for=Enum2
const (
	Enum2Foo = "foo"
)
