package nntp

import (
	"strings"
	"testing"
)

func TestCapabilities(t *testing.T) {
	cli := getTestClient(foo)

	res, err := cli.Capabilities()

	if err != nil {
		t.Fatalf("Bad response for capabilities: %v", err)
	}

	if len(res) != 3 {
		t.Errorf("Wrong number of capabilities in response: %d, %+v", len(res), res)
	}
}

var foo = strings.Replace(`
101 Capabilities List:
FOO
BAR
BAZ
.
`, "\n", "\r\n", -1)
