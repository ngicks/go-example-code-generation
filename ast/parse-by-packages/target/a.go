package target

import "io"

var _ io.Reader = (*Foo)(nil)

type Foo struct{}

func (f *Foo) Read(p []byte) (int, error) {
	return copy(p, []byte(`foo`)), nil
}
