package nntp

import "io"

type Response struct {
	Code    int                 `json:"code"xml:"code"`
	Message string              `json:"message"xml:"message"`
	Headers map[string][]string `json:"headers"xml:"headers"`
	Body    io.ReadCloser       `json:"body"xml:"body"`
}
