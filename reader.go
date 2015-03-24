package nntp

import (
	"bufio"
	"bytes"
	"io"
	"sync"
)

const EndLine = ".\r\n"

//body is an io.Reader that will
type body struct {
	br   *bufio.Reader
	done *sync.WaitGroup
	eof  bool
}

func (b *body) Read(p []byte) (int, error) {
	written := 0

	if b.eof {
		return written, io.EOF
	}

	var bs []byte
	var err error

	for written < len(p) {
		bs, err = b.br.ReadBytes('\n')

		if bytes.Equal(bs, []byte(EndLine)) {
			b.eof = true
			b.done.Done()
			return written, io.EOF
		}

		if len(bs) > 2 && bs[len(bs)-2] == '\r' {
			bs[len(bs)-2] = '\n'
			bs = bs[:len(bs)-1]
		}

		n := copy(p[written:], bs)
		written += n

		if err != nil {
			break
		}
	}

	return written, err
}

func NewArticleReader(r io.Reader) io.Reader {
	var br *bufio.Reader

	switch r := r.(type) {
	case *body:
		return r
	case *bufio.Reader:
		br = r
	default:
		br = bufio.NewReader(r)
	}

	b := body{br, &sync.WaitGroup{}, false}
	b.done.Add(1)
	return &b
}
