package nntp

import (
	"strings"
	"testing"
)

func TestList(t *testing.T) {
	c := getTestClient(listString)

	groups, err := c.List()

	if err != nil {
		t.Fatalf("Error getting groups: %v", err)
	}

	if len(groups) != 2 {
		t.Fatalf("Didn't get the right number of groups back.")
	}

	if groups[0].Id != "alt.foo.bar" {
		t.Errorf("Wrong id for first group.")
	}
}

var listString = strings.Replace(`
215 NewsGroups Follow
alt.foo.bar 78 2 y
alt.bar.bang 123 2 m
.
`, "\n", "\r\n", -1)
