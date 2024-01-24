package nntp

import (
	"strconv"
	"strings"
	"testing"
)

func TestListGroup(t *testing.T) {
	cli := getTestClient(listGroupText)

	arts, err := cli.ListGroup("123")

	if err != nil {
		t.Errorf("Error listing group: %v", err)
	}

	if len(arts) != 5 {
		t.Errorf("Wrong number of articles: %d", len(arts))
	}

	artExp := []int{2, 3, 4, 5, 6}

	for i := range arts {
		if arts[i] != strconv.Itoa(artExp[i]) {
			t.Errorf("Wrong article at index %d: %s, should be %d", i, arts[i], artExp[i])
		}
	}
}

var listGroupText = strings.Replace(`
211 group listings follow
2
3
4
5
6
.
`, "\n", "\r\n", -1)
