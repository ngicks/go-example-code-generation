package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
)

const src = `// test example package
package example

import (
	"fmt"
	unused "bytes"
)

func brokenIndent() {
			// triple indent
	var (
	i int
			j bool
)

	fmt.Printf("i=%d,j=%t\n",i,j) // no space after comma
}
`

func checkGoimports() error {
	_, err := exec.LookPath("goimports")
	return err
}

type cmdPipeReader struct {
	cmd      *exec.Cmd
	pipe     io.Reader
	stderr   *bytes.Buffer
	waitOnce sync.Once
	err      atomic.Value
}

func newCmdPipeReader(cmd *exec.Cmd, r io.Reader) (*cmdPipeReader, error) {
	stderr := new(bytes.Buffer)

	cmd.Stdin = r
	cmd.Stderr = stderr

	p, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return &cmdPipeReader{cmd: cmd, pipe: p, stderr: stderr}, nil
}

func (r *cmdPipeReader) Read(p []byte) (n int, err error) {
	var firstRead bool
	r.waitOnce.Do(func() {
		firstRead = true
		go func() {
			err := r.cmd.Wait()
			if err == nil {
				err = io.EOF
			} else {
				err = fmt.Errorf("%s failed: err = %w, msg = %s", r.cmd.Path, err, r.stderr.Bytes())
			}
			r.err.Store(err)
		}()
	})

	if !firstRead {
		errValue := r.err.Load()
		if errValue != nil {
			return n, errValue.(error)
		}
	}

	n, err = r.pipe.Read(p)
	if err != io.EOF && !errors.Is(err, os.ErrClosed) {
		return n, err
	}

	return n, nil
}

func applyGoimportsPiped(ctx context.Context, r io.Reader) (io.Reader, error) {
	cmd := exec.CommandContext(ctx, "goimports")
	return newCmdPipeReader(cmd, r)
}

func applyGoimports(ctx context.Context, r io.Reader) (*bytes.Buffer, error) {
	r, err := applyGoimportsPiped(ctx, r)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)

	return &buf, err
}

func main() {
	err := checkGoimports()
	if err != nil {
		panic(err)
	}
	fmt.Println("found goimports")

	r, err := applyGoimportsPiped(context.Background(), strings.NewReader(src))
	if err != nil {
		panic(err)
	}
	fmt.Println("success:")
	readAllPrint(r)
	fmt.Println()

	buf, err := applyGoimports(context.Background(), strings.NewReader(src))
	if err != nil {
		panic(err)
	}
	fmt.Printf("formatted: \n---\n%s\n---\n", buf.Bytes())

	r, err = applyGoimportsPiped(context.Background(), strings.NewReader(src[:len(src)-5]))
	if err != nil {
		panic(err)
	}
	fmt.Println("error:")
	readAllPrint(r)
}

func readAllPrint(r io.Reader) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("formatted: \n---\n%s\n---\n", buf.Bytes())
}
