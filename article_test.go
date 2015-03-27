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

func TestArticleGetter(t *testing.T) {
	d := &Client{
		Connections: 1,
		Username:    "andrew",
		Password:    "1234",
		pool:        make(chan *connection, 1),
	}

	buf := &bytes.Buffer{}

	d.pool <- newConnection(testrw{
		Writer: buf,
		Reader: strings.NewReader(testServer2),
	})

	r, err := d.GetArticle("foo")

	if err != nil {
		t.Fatal("failed to get article: %v", err)
	}

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

	if s != "test\r\n" {
		t.Errorf("Wrong string read -> %s", s)
	}

	if s, err = bufr.ReadString('\n'); err != nil {
		t.Errorf("String was not read: %v", err)
	} else if s != ".foo\r\n" {
		t.Errorf("Wrong string read: len: %d, %s", len(s), s)
	}

	s, err = bufr.ReadString('\n')

	if err != io.EOF {
		t.Errorf("Missed eof: %s", s)
	}
}

var testServer2 = strings.Replace(`220 found
Foo: Bar

test
..foo
`, "\n", "\r\n", -1)
