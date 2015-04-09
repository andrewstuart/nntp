package nntp

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func getTestClient(s string) *Client {
	cli := &Client{
		p: &sync.Pool{},
	}

	tc := &testCloser{
		Reader: strings.NewReader(s),
	}
	cli.p.Put(NewConn(tc))
	return cli
}

func TestClient(t *testing.T) {
	cli := getTestClient(testClientResponse)

	res, err := cli.Do("FOO BAR")

	if err != nil {
		t.Fatalf("Error doing FOO BAR: %v", err)
	}

	if res.Code != 220 {
		t.Errorf("Wrong code")
	}

	if len(res.Headers) != 2 {
		t.Errorf("Wrong number of headers: %d", len(res.Headers))
	}

	buf := &bytes.Buffer{}

	if res.Body == nil {
		t.Fatalf("Did not return a body")
	}

	io.Copy(buf, res.Body)

	if buf.String() != "foobarbaz\r\n.foo\r\n" {
		t.Errorf("Wrong body returned: %s", buf.String())
	}

	if cli.p.Get() != nil {
		t.Errorf("Somehow got non-nil from pool")
	}

	err = res.Body.Close()

	if err != nil {
		t.Errorf("Error closing body: %v", err)
	}

	if cli.p.Get() == nil {
		t.Errorf("Could not get a connection after closing body")
	}
}

func TestNoBody(t *testing.T) {
	cli := getTestClient(testClientNoBody)

	res, err := cli.Do("FOO BAR")

	if err != nil {
		t.Fatal(err)
	}

	if res.Body != nil {
		t.Errorf("Got a body when res code indicated no body")
	}

	if res.Code != 500 {
		t.Errorf("Wrong code: %d", res.Code)
	}

	if conn := cli.p.Get(); conn == nil {
		t.Fatalf("Could not get a connection from the pool -- should have been replaced.")
	}
}

func TestS(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, a)
	}))
	defer ts.Close()
}

var testClientNoBody = "500 No Body\r\n"

var testClientResponse = strings.Replace(`220 Have Body
H: bar
B: foo

foobarbaz
..foo
.
`, "\n", "\r\n", -1)
