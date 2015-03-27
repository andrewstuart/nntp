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
		pool:        make(chan *connection, 1),
	}

	c := newConnection(trw2)

	cli.pool <- c

	//TODO use res
	_, err := cli.Do("FOO")

	if err != nil {
		t.Errorf("message")
	}

}
