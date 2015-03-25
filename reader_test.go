package nntp

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	r := NewReader(strings.NewReader(readerTest))

	bs, err := ioutil.ReadAll(r)

	if err != nil {
		t.Fatalf("error while reading: %v")
	}

	if bytes.Contains(bs, []byte("..Whodunit")) {
		t.Errorf("Double dot should have been escaped")
	}
}

var readerTest = strings.Replace(`
Header1: Foo
Header2: Bar
Header3: Baz

..Whodunit
This is the question
I don't really know
.`, "\n", "\r\n", -1)
