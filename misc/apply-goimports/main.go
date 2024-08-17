package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
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
	err      error
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
	return r.pipe.Read(p)
}

func (r *cmdPipeReader) Close() error {
	r.waitOnce.Do(func() {
		err := r.cmd.Wait()
		if err != nil {
			err = fmt.Errorf("%s failed: err = %w, msg = %s", r.cmd.Path, err, r.stderr.Bytes())
		}
		r.err = err
	})
	return r.err
}

func applyGoimportsPiped(ctx context.Context, r io.Reader) (io.ReadCloser, error) {
	cmd := exec.CommandContext(ctx, "goimports")
	return newCmdPipeReader(cmd, r)
}

func applyGoimports(ctx context.Context, r io.Reader) (*bytes.Buffer, error) {
	p, err := applyGoimportsPiped(ctx, r)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, p)
	cErr := p.Close()

	switch {
	case err != nil && cErr != nil:
		err = fmt.Errorf("copy err: %w, wait err: %w", err, cErr)
	case err != nil:
	case cErr != nil:
		err = cErr
	}

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

func readAllPrint(r io.ReadCloser) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	cErr := r.Close()
	if err != nil || cErr != nil {
		fmt.Printf("copy err: %v, wait err: %v\n", err, cErr)
		return
	}
	fmt.Printf("formatted: \n---\n%s\n---\n", buf.Bytes())
}
