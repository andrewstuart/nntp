package nntp

import (
	"bufio"
	"bytes"
	"fmt"
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
	if b.eof {
		return 0, io.EOF
	}

	if b.br == nil {
		return 0, io.ErrClosedPipe
	}

	//Read until newline
	bs, err := b.br.ReadBytes('\n')

	if err != nil {
		return 0, fmt.Errorf("error reading bytes: %v", err)
	}

	//Check EOF
	if bytes.Equal(bs, []byte(EndLine)) {
		b.done.Done()
		b.eof = true
		return 0, io.EOF
	}

	//Remove silly carriage returns
	if len(bs) > 2 && bs[len(bs)-2] == '\r' {
		//Remove last byte
		bs = bs[:len(bs)-1]
		//Change new last byte to \n
		bs[len(bs)-1] = '\n'
	}

	//Drop leading first byte if needed
	if len(bs) > 2 && bytes.Equal(bs, []byte("..")) {
		bs = bs[1:]
	}

	//Copy and return
	return copy(p, bs), nil
}

func NewArticleReader(r io.Reader) io.Reader {
	var br *bufio.Reader
	switch r := r.(type) {
	case *body:
		return r
	case *bufio.Reader:
		//*connection goes here??
		br = r
	default:
		br = bufio.NewReader(r)
	}

	b := body{br, &sync.WaitGroup{}, false}
	b.done.Add(1)
	return &b
}
