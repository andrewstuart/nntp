package nntp

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type Conn struct {
	*Reader
	io.Writer
	cls chan bool
}

func NewConn(c net.Conn) io.ReadWriteCloser {
	return &Conn{
		Reader: NewReader(bufio.NewReader(c)),
		Writer: io.Writer(c),
	}
}

func (c *Conn) Do(format string, is ...interface{}) *Response {
	cmd := fmt.Sprintf(format, is...)
	fmt.Fprintf(c, "%s\r\n", cmd)
	return nil
}
