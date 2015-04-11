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

func NewConn(c io.ReadWriteCloser) *Conn {
	br := bufio.NewReader(c)
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

	res, err := NewResponse(c.br)

	fmt.Printf(format+"\n", is...)

	if err != nil {
		return nil, err
	}

	return res, nil
}
