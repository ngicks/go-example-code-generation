package target

//enum:variants=foo,bar,baz
type Enum string

//enum:generated_for=Enum
const (
	EnumFoo	Enum	= "foo"
	EnumBar	Enum	= "bar"
	EnumBaz	Enum	= "baz"
)

//enum:variants=foo,bar,baz
type Enum2 string

//enum:generated_for=Enum2
const (
	Enum2Foo	Enum2	= "foo"
	Enum2Bar	Enum2	= "bar"
	Enum2Baz	Enum2	= "baz"
)
