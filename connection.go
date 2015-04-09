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
	*bufio.Reader
	io.Writer
	cls chan (chan error)
}

func NewConn(c io.ReadWriteCloser) *Conn {
	return &Conn{
		Reader: bufio.NewReader(c),
		Writer: io.Writer(c),
		cls:    make(chan (chan error)),
	}
}

func (c *Conn) Close() error {
	return nil
}

func (c *Conn) Do(format string, is ...interface{}) (*Response, error) {
	fmt.Fprintf(c, strings.TrimSpace(format)+"\r\n", is...)
	res, err := NewResponse(c)

	if err != nil {
		return nil, err
	}

	return res, nil
}
