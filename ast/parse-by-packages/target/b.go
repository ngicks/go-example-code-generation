package target

import "fmt"

func Bar() error {
	p := make([]byte, len("foo"))
	_, err := new(Foo).Read(p)
	if err != nil {
		return err
	}
	fmt.Printf("foo:%s\n", p)
	return nil
}
