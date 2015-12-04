package nntp

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	r := NewReader(bufio.NewReader(strings.NewReader(readerTest)))

	bs, err := ioutil.ReadAll(r)

	if err != nil {
		t.Fatalf("error while reading: %v", err)
	}

	lines := strings.Split(string(bs), "\n")

	if len(lines) != 9 {
		t.Errorf("Wrong number of lines: %d", len(lines))
	}

	dotLine := 5
	if lines[dotLine] != ".Whodunit\r" {
		t.Errorf("Did not properly escape double dot: %s", lines[dotLine])
	}
}

func BenchmarkReader(b *testing.B) {
	b.SetBytes(int64(len(readerTest)))
	b.RunParallel(func(pb *testing.PB) {
		bs := []byte(readerTest)
		dest := &bytes.Buffer{}

		for pb.Next() {
			r := NewReader(bytes.NewReader(bs))

			n, err := io.Copy(dest, r)

			if err != nil {
				b.Fatalf("error while reading: %v", err)
			}

			if n == 0 {
				b.Fatalf("Read zero bytes")
			}
		}
	})
}

var readerTest = strings.Replace(`201 found
Header1: Foo
Header2: Bar
Header3: Baz

..Whodunit
This is the question
I don't really know
.`, "\n", "\r\n", -1)
