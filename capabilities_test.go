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

	_, err = cli.Capabilities()

	if err == nil {
		t.Errorf("Did not return error for bad code: %v", err)
	}
}

var foo = strings.Replace(`
101 Capabilities List:
FOO
BAR
BAZ
.
502 Capabilites are dumb
.
`, "\n", "\r\n", -1)
