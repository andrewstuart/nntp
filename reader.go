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

func (r *Reader) Read(p []byte) (written int, err error) {
	if r.eof {
		err = io.EOF
	}

	var bt byte
	for err == nil && written < len(p) {
		bt, err = r.br.ReadByte()

		if err != nil {
			return
		}

		switch bt {
		case '.':
			if r.nl {
				var bs []byte
				bs, err = r.br.Peek(2)

				if err != nil {
					return
				}

				if bytes.Equal(bs, []byte("\r\n")) || bs[0] == '\n' {
					r.eof = true
					err = io.EOF
					r.br.ReadBytes('\n')
					return
				} else {
					r.nl = false
					continue
				}
			}
		case '\n':
			r.nl = true
		}

		p[written] = bt
		written++

		if bt != '\n' {
			r.nl = false
		}
	}
	return
}
