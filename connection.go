package nntp

import (
	"bufio"
	"fmt"
	"io"
)

type ConnErr struct {
	Code   int    `json:"code"xml:"code"`
	Reason string `json:"reason"xml:"reason"`
}

func (c ConnErr) Error() string {
	return fmt.Sprintf("%d: %s", c.Code, c.Reason)
}

type Conn struct {
	*Reader
	io.Writer
	onClose func() error
	cls     chan (chan error)
}

func NewConn(c io.ReadWriteCloser) *Conn {
	return &Conn{
		Reader: NewReader(bufio.NewReader(c)),
		Writer: io.Writer(c),
		cls:    make(chan (chan error)),
	}
}

func (c *Conn) Close() error {
	if c.onClose != nil {
		return c.onClose()
	}
	return nil
}

func (c *Conn) Do(format string, is ...interface{}) (*Response, error) {
	fmt.Fprintf(c, format+"\r\n", is...)
	res, err := NewResponse(c.R)

	if err != nil {
		return nil, err
	}

	return res, nil
}
