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

type Response struct {
	Code    int                  `json:"code"xml:"code"`
	Message string               `json:"message"xml:"message"`
	Headers textproto.MIMEHeader `json:"headers"xml:"headers"`
	Body    io.ReadCloser        `json:"body"xml:"body"` //Presence (non-nil) indicates multi-line response
	br      *bufio.Reader
}

var isMultiLine = map[int]bool{
	100: true,
	101: true,
	211: true,
	215: true,
	220: true,
	221: true,
	222: true,
	224: true,
	225: true,
	230: true,
	231: true,
}

func NewResponse(r io.Reader) (*Response, error) {
	nr := NewReader(r)

	s, err := nr.R.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading header: %v", err)
	}

	sa := strings.Split(strings.TrimSpace(s), " ")
	if len(sa) < 2 {
		return nil, IllegalResponse
	}

	res := &Response{
		br:      nr.R,
		Headers: make(map[string][]string),
	}

	if code, err := strconv.Atoi(sa[0]); err != nil {
		return nil, IllegalResponse
	} else {
		res.Code = code
		res.Message = sa[1]
	}

	tpr := textproto.NewReader(nr.R)

	if isMultiLine[res.Code] {
		h, err := tpr.ReadMIMEHeader()

		if err != nil {
			return nil, IllegalResponse
		}

		res.Headers = h
		res.Body = nr
	}

	return res, nil
}
