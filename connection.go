package nntp

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type ConnErr struct {
	Code   int    `json:"code"xml:"code"`
	Reason string `json:"reason"xml:"reason"`
}

func (c ConnErr) Error() string {
	return fmt.Sprintf("%d: %s", c.Code, c.Reason)
}

type Conn struct {
	br *bufio.Reader
	w  io.Writer
}

func (c *Conn) Wrap(fn ...func(io.Reader) io.Reader) error {
	var r io.Reader
	for i := range fn {
		r = fn[i](r)
	}
	c.br = bufio.NewReader(r)
	return nil
}

func NewConn(c io.ReadWriteCloser, wrappers ...func(io.Reader) io.Reader) *Conn {
	var r io.Reader
	for w := range wrappers {
		r = wrappers[w](c)
	}

	br := bufio.NewReader(r)
	nnConn := Conn{
		br: br,
		w:  c,
	}

	//Throw away welcome line
	br.ReadBytes('\n')

	return &nnConn
}

func (c *Conn) Close() error {
	return nil
}

func (c *Conn) Do(format string, is ...interface{}) (*Response, error) {
	fmt.Fprintf(c.w, strings.TrimSpace(format)+"\r\n", is...)
	return NewResponse(c.br)
}
