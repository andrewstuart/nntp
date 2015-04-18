package nntp

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"testing"

	"git.astuart.co/andrew/pool"
)

func getTestClient(s string) *Client {
	cli := &Client{
		p: pool.NewPool(nil),
	}

	tc := &testCloser{
		Reader: strings.NewReader(s),
	}

	cli.p.Put(NewConn(tc))
	return cli
}

func TestClient(t *testing.T) {
	go http.ListenAndServe(":6060", nil)
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

	err = res.Body.Close()

	if err != nil {
		t.Errorf("Error closing body: %v", err)
	}

	c, _ := cli.p.Get()
	if c.(*Conn) == nil {
		t.Errorf("Could not get a connection after closing body")
	}
	cli.p.Put(c)

	res, err = cli.Do("Foo again")

	if err != nil {
		t.Fatalf("Second read failed for test client")
	}

	if res.Code != 200 {
		t.Errorf("Wrong response code: %d", res.Code)
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

	if conn, _ := cli.p.Get(); conn == nil {
		t.Fatalf("Could not get a connection from the pool -- should have been replaced.")
	}
}

var testClientNoBody = "Welome\r\n500 No Body\r\n"

var testClientResponse = strings.Replace(`Fooba Welcome
220 Have Body
H: bar
B: foo

foobarbaz
..foo
.
200 no body
`, "\n", "\r\n", -1)

func TestNewClient(t *testing.T) {
	go func() {
		ln, err := net.Listen("tcp", ":15531")

		if err != nil {
			t.Errorf("Error setting up test conn: %v", err)
		}

		conn, err := ln.Accept()

		if err != nil {
			t.Fatalf("Could not accept. %v", err)
		}
		fmt.Fprint(conn, testClientResponse)
		conn.Close()
	}()

	cli := NewClient("localhost", 15531)
	cli.SetMaxConns(1)

	res, err := cli.Do("foo bar")

	if err != nil {
		t.Errorf("Error on NewClient test: %v", err)
	}

	if res.Body == nil {
		t.Fatalf("Did not have a body on response")
	}

	buf := &bytes.Buffer{}
	io.Copy(buf, res.Body)

	res.Body.Close()

	if len(res.Headers) != 2 {
		t.Errorf("Wrong number of headers: %d", len(res.Headers))
	}

	if res.Code != 220 {
		t.Errorf("Wrong res code on NewClient: %d", res.Code)
	}

	if buf.String() != "foobarbaz\r\n.foo\r\n" {
		t.Errorf("Wrong body returned: %s", buf)
	}
}
