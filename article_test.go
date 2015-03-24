package nntp

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"
)

type testrw struct {
	io.Writer
	io.Reader
}

var trw = testrw{
	&bytes.Buffer{},
	strings.NewReader("220 found\r\nFoo: Bar\r\n\r\ntest\r\n.\r\n"),
}

func TestArticleGetter(t *testing.T) {
	d := &Client{
		Connections: 1,
		Username:    "andrew",
		Password:    "1234",
		cBucket:     make(chan *connection, 1),
	}

	d.cBucket <- newConnection(trw)

	r, err := d.GetArticle("foo")

	if r.Headers["Foo"] != "Bar" {
		t.Errorf("Headers were read improperly")
	}

	if err != nil {
		t.Errorf("got an error for GetArticle: %v", err)
	}

	bufr := bufio.NewReader(r.Body)

	s, err := bufr.ReadString('\n')

	if err != nil {
		t.Fatalf("error reading string: %v", err)
	}

	if strings.TrimSpace(s) != "test" {
		t.Errorf("Wrong string read")
	}

	s, err = bufr.ReadString('\n')

	if err != io.EOF {
		t.Errorf("Did not read EOF properly")
	}
}
