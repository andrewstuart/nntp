package nntp

import (
	"bufio"
	"bytes"
	"io"
)

const EndLine = ".\r\n"

var EndBytes = []byte(EndLine)

type Reader struct {
	br      *bufio.Reader
	eof, nl bool
}

func NewReader(r io.Reader) *Reader {
	switch r := r.(type) {
	case *Reader:
		return r
	default:
		return &Reader{
			br: bufio.NewReader(r),
		}
	}
}

func (b *Reader) Read(p []byte) (written int, err error) {
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
