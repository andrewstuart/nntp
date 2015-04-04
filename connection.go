package nntp

import (
	"fmt"
	"io"
)

type Conn struct {
}

func (c *Conn) Write(p []byte) (int, error) {
	return 0, nil
}

func (c *Conn) Read(p []byte) (n int, err error) {
	return
}

func (c *Conn) do(format string, is ...interface{}) io.ReadCloser {
	cmd := fmt.Sprintf(format, is...)
	fmt.Fprintf(c, "%s\r\n", cmd)
	return nil
}
