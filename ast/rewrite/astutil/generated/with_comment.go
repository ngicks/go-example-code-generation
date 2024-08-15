package target

// free floating comment 1

func Foo() {
	// nothing
}

//enum:variants=foo,bar,baz,qux,quux,corge
type EnumWithComments string

// free floating comment 2
//enum:generated_for=EnumWithComments

// nothing
const (
	EnumWithCommentsFoo	EnumWithComments	= "foo"
	EnumWithCommentsBar	EnumWithComments	= "bar"
	EnumWithCommentsBaz	EnumWithComments	= "baz"
	EnumWithCommentsQux	EnumWithComments	= "qux"
	EnumWithCommentsQuux	EnumWithComments	=

	//enum:variants=foo,bar,baz
	"quux"
	EnumWithCommentsCorge	EnumWithComments	= "corge"
)

func Bar() {

}

type EnumWithComments2 string

//enum:generated_for=EnumWithComments2
const (
	EnumWithComments2Foo	EnumWithComments2	= "foo"
	EnumWithComments2Bar	EnumWithComments2	= "bar"
	EnumWithComments2Baz	EnumWithComments2	= "baz"
)

/* free floating comment 4


 */
