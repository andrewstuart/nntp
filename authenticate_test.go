package nntp

import "testing"

func TestAuth(t *testing.T) {
	cli := getTestClient("381 PASS\r\n281 Accept\r\n")
	err := cli.Auth("foo", "bar")
	if err != nil {
		t.Errorf("Error authenticating: %v", err)
	}
	if cli.p.Get() == nil {
		t.Errorf("Couldn't get back a connection")
	}
}

func TestAuth2(t *testing.T) {
	cli := getTestClient("281 Accept\r\n")
	err := cli.Auth("foo", "")
	if err != nil {
		t.Errorf("Error authenticating: %v", err)
	}
}
