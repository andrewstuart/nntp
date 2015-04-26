package nntp

import (
	"strings"
	"testing"
)

func TestHead(t *testing.T) {
	tc := getTestClient(headTest)

	res, err := tc.Do("HEAD 1")

	if err != nil {
		t.Errorf("Got an error on head request: %v", err)
	}

	if len(res.Headers) != 2 {
		t.Errorf("wrong number of headers: %d", len(res.Headers))
	}
}

var headTest = strings.Replace(`
221 Yup
ABC: here they are
DEF: foo 2 3
.
`, "\r", "\r\n", -1)
