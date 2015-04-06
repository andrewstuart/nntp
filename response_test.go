package nntp

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestResponse(t *testing.T) {
	resReader := strings.NewReader(resString)

	res, err := NewResponse(resReader)

	if err != nil {
		t.Errorf("error getting test response: %v", err)
	}

	if res.Code != 201 {
		t.Errorf("Wrong response code: %d, should be 201", res.Code)
	}

	if res.Message != "Foo" {
		t.Errorf("Wrong test response code: %s, should be %s", res.Message, "Foo")
	}

	if len(res.Headers) < 2 {
		t.Errorf("Wrong number of headers: %d, should be 2", len(res.Headers))
	}

	a := map[string]string{
		"Header":     "1",
		"Header-Two": "2",
	}

	for k, v := range a {
		h := res.Headers[k]
		if len(h) < 1 {
			t.Fatalf("Wrong number of headers for key %s: %d", k, len(h))
		}

		if h[0] != v {
			t.Errorf("Wrong header returned for key %s: %s", k, v)
		}
	}

	bs, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Errorf("error reading body: %v", err)
	}

	strAr := strings.Split(string(bs), "\n")

	if len(strAr) < 4 {
		t.Fatalf("Wrong number of lines in Body")
	}

	if strAr[1] != ".Foo man chu\r" {
		t.Errorf("Wrong body: %s", strAr[1])
	}
}

var resString = strings.Replace(`201 Foo
Header: 1
Header-Two: 2

Whatever
..Foo man chu
I like this stuff
.
`, "\n", "\r\n", -1)
