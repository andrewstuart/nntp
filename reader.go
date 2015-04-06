package nntp

import (
	"bufio"
	"bytes"
	"io"
)

const EndLine = ".\r\n"

var EndBytes = []byte(EndLine)

type Reader struct {
	R       *bufio.Reader
	eof, nl bool
	c       io.Closer //We'll call close if possible on underlying reader
}

func NewReader(r io.Reader) *Reader {
	switch r := r.(type) {
	case *Reader:
		return r
	default:
		rr := Reader{
			R: bufio.NewReader(r),
		}

		if c, isCloser := r.(io.Closer); isCloser {
			rr.c = c
		}

		return &rr
	}
}

func (r *Reader) Read(p []byte) (written int, err error) {
	if r.eof {
		err = io.EOF
	}

	var bt byte
	for err == nil && written < len(p) {
		bt, err = r.R.ReadByte()

		if err != nil {
			return
		}

		switch bt {
		case '.':
			if r.nl {
				var bs []byte
				bs, err = r.R.Peek(2)

				if err != nil {
					return
				}

				if bytes.Equal(bs, []byte("\r\n")) || bs[0] == '\n' {
					r.eof = true
					err = io.EOF
					r.R.ReadBytes('\n')
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

func (r *Reader) Close() error {
	if r.c != nil {
		return r.c.Close()
	}

	return nil
}

//Next enables the use of io.Readers that may have multiple bodies.
func (r *Reader) Next() (*Reader, error) {
	if _, err := r.R.Peek(1); err != nil {
		return nil, err
	}

	r.eof = false
	r.nl = true

	return r, nil
}
