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

func (b *body) Read(p []byte) (int, error) {
	written := 0
	var err error

	if b.eof {
		err = io.EOF
	}

	var bt byte
readLoop:
	for err == nil && written < len(p) {
		bt, err = b.br.ReadByte()

		if err != nil {
			break readLoop
		}

		switch bt {
		case '.':
			if b.nl {
				var bs []byte
				bs, err = b.br.Peek(2)

				if err != nil {
					break readLoop
				}

				if bytes.Equal(bs, []byte("\r\n")) {
					b.eof = true
					b.done.Done()
					err = io.EOF
					b.br.ReadBytes('\n')
					break readLoop
				} else {
					b.nl = false
					continue
				}
			}
		case '\n':
			b.nl = true
		case '\r':
			continue readLoop
		}

		p[written] = bt
		written++

		if bt != '\n' {
			b.nl = false
		}
	}
	return written, err
}

func NewReader(r io.Reader) io.Reader {
	var br *bufio.Reader

	if r, isBody := r.(*body); isBody {
		return r
	}

	br = bufio.NewReader(r)

	b := body{
		br:   br,
		done: &sync.WaitGroup{},
	}
	b.done.Add(1)
	return &b
}
