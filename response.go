package nntp

import (
	"bufio"
	"bytes"
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
	br := bufio.NewReader(r)

	s, err := br.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading response line: %v", err)
	}

	sa := strings.Split(strings.TrimSpace(s), " ")
	if len(sa) < 2 {
		return nil, fmt.Errorf("error getting response code: %v", err)
	}

	res := &Response{
		br:      br,
		Headers: make(map[string][]string),
	}

	if code, err := strconv.Atoi(sa[0]); err != nil {
		return nil, fmt.Errorf("error converting response code: %v", err)
	} else {
		res.Code = code
		res.Message = strings.Join(sa[1:], " ")
	}

	if isMultiLine[res.Code] || res.Code == 211 && strings.Contains(res.Message, "follow") {
		switch res.Code {
		case 221, 222, 220:
			h, err := readHeaders(br)
			res.Headers = h

			//Early return if EOF (no Body)
			if err == io.EOF {
				return res, nil
			}
		}

		res.Body = NewReader(br)
	}

	return res, nil
}

func readHeaders(br *bufio.Reader) (map[string][]string, error) {
	h := make(map[string][]string)
	var err error
	var bs []byte

	for err == nil {
		bs, err = br.ReadBytes('\n')

		isEOF := bytes.Equal(bs, []byte(".\r\n"))
		isEnd := bytes.Equal(bs, []byte("\r\n"))
		if err == io.EOF || isEOF {
			return h, io.EOF
		} else if isEnd {
			return h, nil
		}

		kv := bytes.Split(bytes.TrimSpace(bs), []byte(": "))

		if len(kv) < 2 {
			return h, fmt.Errorf("malformed header")
		}

		k := string(kv[0])
		v := string(kv[1])

		if _, ok := h[k]; !ok {
			h[k] = make([]string, 0, 1)
		}
		h[k] = append(h[k], v)
	}

	return h, err
}
