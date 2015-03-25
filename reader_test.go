package nntp

import (
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

	t.Log(string(bs))

}

var readerTest = `
Header1: Foo
Header2: Bar
Header3: Baz

..Whodunit
This is the question
I don't really know
.` + "\r\n"
