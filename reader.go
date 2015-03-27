package nntp

import (
	"bufio"
	"bytes"
	"io"
	"sync"
)

const EndLine = ".\r\n"

var EndBytes = []byte(EndLine)

//body is an io.Reader that will
type body struct {
	br      *bufio.Reader
	done    *sync.WaitGroup
	eof, nl bool
}

func newBody(r io.Reader) *body {
	if r, isBody := r.(*body); isBody {
		return r
	}

	b := body{
		br:   bufio.NewReader(r),
		done: &sync.WaitGroup{},
	}
	b.done.Add(1)
	return &b
}

func (b *body) Read(p []byte) (written int, err error) {
	if b.eof {
		err = io.EOF
	}

	var bt byte
	for err == nil && written < len(p) {
		bt, err = b.br.ReadByte()

		if err != nil {
			return
		}

		switch bt {
		case '.':
			if b.nl {
				var bs []byte
				bs, err = b.br.Peek(2)

				if err != nil {
					return
				}

				if bytes.Equal(bs, []byte("\r\n")) || bs[0] == '\n' {
					b.eof = true
					b.done.Done()
					err = io.EOF
					b.br.ReadBytes('\n')
					return
				} else {
					b.nl = false
					continue
				}
			}
		case '\n':
			b.nl = true
		}

		p[written] = bt
		written++

		if bt != '\n' {
			b.nl = false
		}
	}
	return
}

func NewReader(r io.Reader) io.Reader {
	return newBody(r)
}
