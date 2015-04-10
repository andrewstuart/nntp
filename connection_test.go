package nntp

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"
)

type testCloser struct {
	io.Reader
	closed  bool
	written int
}

func (tc *testCloser) Close() error {
	tc.closed = true
	return nil
}

func (tc *testCloser) Write(p []byte) (int, error) {
	tc.written += len(p)
	return len(p), nil
}

func TestConnection(t *testing.T) {
	tc := &testCloser{
		Reader: bufio.NewReader(strings.NewReader(clientTestString)),
	}
	nc := NewConn(tc)

	res, err := nc.Do("SOME THING")

	if err != nil {
		t.Fatalf("Error on client Do: %v", err)
	}

	b := &bytes.Buffer{}

	if res.Body == nil {
		t.Fatalf("Response should have contained a body")
	}

	_, err = io.Copy(b, res.Body)

	if err != nil {
		t.Errorf("Error copying body: %v", err)
	}

	if b.String() != "Foo\r\n" {
		t.Errorf("Wrong body reported: %s", b.String())
	}

	err = res.Body.Close()

	if err != nil {
		t.Errorf("Error closing body: %v", err)
	}

	if tc.closed {
		t.Errorf("Body Close closed underlying test connection")
	}

	_, err = nc.Do("foo")

	if err == nil {
		t.Errorf("Error was not returned on empty reader")
	}

	err = nc.Close()

}

var clientTestString string = strings.Replace(`220 okay
H: FooBar

Foo
.
`, "\n", "\r\n", -1)
