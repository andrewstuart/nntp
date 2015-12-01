package nntp

import (
	"bufio"
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

var nl = []byte{'\r', '\n'}

//The Read method handles translation of the NNTP escaping and marking EOF when
//the end of a body is received.
func (r *Reader) Read(p []byte) (bytesRead int, err error) {
	if r.eof {
		err = io.EOF
		return
	}

	var bt byte
	var bs []byte
	for err == nil && bytesRead < len(p) {
		bt, err = r.R.ReadByte()

		if err != nil {
			return
		}

		switch bt {
		case '.':
			if r.nl {
				bs, err = r.R.Peek(2)

				if err != nil {
					return
				}

				if len(bs) == 2 && bs[0] == '\r' && bs[1] == '\n' {
					r.eof = true
					err = io.EOF
					r.R.ReadBytes('\n')
					return
				}

				r.nl = false
				continue
			}
		case '\n':
			r.nl = true
		}

		p[bytesRead] = bt
		bytesRead++

		if bt != '\n' {
			r.nl = false
		}
	}
	return
}

//Close allows users of a Reader to signal that they are done using the reader.
func (r *Reader) Close() error {
	if r.c != nil {
		//TODO what should I do with this error?
		r.c.Close()
	}

	//Reset reader
	r.eof = false
	r.nl = false

	return nil
}
