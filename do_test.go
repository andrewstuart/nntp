package nntp

import (
	"bytes"
	"strings"
	"testing"
)

var trw2 = testrw{
	&bytes.Buffer{},
	strings.NewReader("220 found\r\ntest\r\n.\r\n"),
}

func TestDo(t *testing.T) {
	cli := &Client{
		Connections: 1,
		Username:    "andrew",
		Password:    "test",
		cBucket:     make(chan *connection, 1),
	}

	c := newConnection(trw2)

	cli.cBucket <- c

	res, err := cli.do("FOO")

	defer func() {
		cli.cBucket <- res.conn
	}()

	if err != nil {
		t.Errorf("message")
	}

	s, err := res.conn.ReadString('\n')

	if err != nil {
		t.Errorf("M2")
	}

	if strings.TrimSpace(s) != "test" {
		t.Errorf("Wrong string: %s", s)
	}
}
