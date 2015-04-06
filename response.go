package nntp

import (
	"bufio"
	"fmt"
	"io"
	"net/textproto"
	"strconv"
	"strings"
)

var (
	IllegalResponse = fmt.Errorf("illegal response")
	IllegalHeader   = fmt.Errorf("illegal headers")
)

type body struct {
	io.Reader
	c io.Closer
}

func (b *body) Close() error {
	if b.c != nil {
		return b.c.Close()
	}
	return nil
}

type Response struct {
	Code    int                  `json:"code"xml:"code"`
	Message string               `json:"message"xml:"message"`
	Headers textproto.MIMEHeader `json:"headers"xml:"headers"`
	Body    io.ReadCloser        `json:"body"xml:"body"`
	br      *bufio.Reader
}

func NewResponse(r io.Reader) (*Response, error) {
	//TODO is there a better way to make sure underlying reader isn't drained by
	//bufio?
	var br *bufio.Reader
	bdy := &body{}

	//Normalize to *Reader
	switch r := r.(type) {
	case (*Reader):
		br = bufio.NewReader(r)
		bdy.c = r.c
	default:
		br = bufio.NewReader(NewReader(r))
	}

	s, err := br.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading header: %v", err)
	}

	sa := strings.Split(strings.TrimSpace(s), " ")
	if len(sa) < 2 {
		return nil, IllegalResponse
	}

	res := &Response{
		br:      br,
		Headers: make(map[string][]string),
	}

	if code, err := strconv.Atoi(sa[0]); err != nil {
		return nil, IllegalResponse
	} else {
		res.Code = code
		res.Message = sa[1]
	}

	tpr := textproto.NewReader(br)

	h, err := tpr.ReadMIMEHeader()

	if err != nil {
		return nil, IllegalResponse
	}

	res.Headers = h

	bdy.Reader = br
	res.Body = bdy

	return res, nil
}
