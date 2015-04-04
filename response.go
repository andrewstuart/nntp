package nntp

import (
	"fmt"
	"io"
)

type Response struct {
	Code    int                 `json:"code"xml:"code"`
	Message string              `json:"message"xml:"message"`
	Headers map[string][]string `json:"headers"xml:"headers"`
	Body    io.ReadCloser       `json:"body"xml:"body"`
}

func (r *Response) Error() string {
	return fmt.Sprintf("nntp error code %d, server says: %s", r.Code, r.Message)
}
