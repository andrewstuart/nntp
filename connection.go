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

	orig io.ReadWriteCloser
}

func (c *Conn) Wrap(fn ...func(io.Reader) io.Reader) {
	if fn != nil {
		var r io.Reader
		for i := range fn {
			r = fn[i](r)
		}
		c.br = bufio.NewReader(r)
	}
}

func NewConn(c io.ReadWriteCloser, wrappers ...func(io.Reader) io.Reader) (string, *Conn, error) {
	var r io.Reader = c
	for w := range wrappers {
		r = wrappers[w](c)
	}

	br := bufio.NewReader(r)
	nnConn := Conn{
		br: br,
		w:  c,

		orig: c,
	}

	//Throw away welcome line
	w, err := br.ReadBytes('\n')
	if err != nil {
		return "", nil, fmt.Errorf("Error reading welcome line: %s", err)
	}

	return string(w), &nnConn, nil
}

func (c *Conn) Close() error {
	return nil
}

func (c *Conn) Do(format string, is ...interface{}) (*Response, error) {
	fmt.Fprintf(c.w, strings.TrimSpace(format)+"\r\n", is...)
	return NewResponse(c.br)
}
