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

func (b body) Read(p []byte) (int, error) {
	written := 0

	if b.eof {
		return written, io.EOF
	}

	for written < len(p) && b.br.Buffered() > 0 {
		bt, err := b.br.ReadByte()

		if err != nil {
			return written, err
		}

		switch bt {
		case '.':
			if bs, err := b.br.Peek(2); err == nil {
				if bytes.Equal(bs, []byte("\r\n")) {
					b.br.ReadByte()
					b.br.ReadByte()
					b.eof = true
					b.done.Done()
					return written, io.EOF
				} else if bs[0] == '.' {
					b.br.ReadByte()
					//Go back to copying
					break
				}
			}
		case '\r':
			bt, err = b.br.ReadByte()

			if err != nil {
				return written, err
			}
		}

		p[written] = bt
		written++
	}
	return written, nil
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
