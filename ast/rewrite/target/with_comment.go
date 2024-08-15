package target

// free floating comment 1

func Foo() {
	// nothing
}

//enum:variants=foo,bar,baz,qux,quux,corge
type EnumWithComments string

// free floating comment 2

func Bar() {
	// nothing
}

//enum:variants=foo,bar,baz
type EnumWithComments2 string

// free floating comment 3

//enum:generated_for=EnumWithComments2
const (
	EnumWithComments2Foo EnumWithComments2 = "foo"
)

/* free floating comment 4


 */
