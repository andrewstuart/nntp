package nntp

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func (cli *Client) Capabilities() ([]string, error) {
	c, err := cli.p.Get()

	if err != nil {
		return nil, fmt.Errorf("error getting pool connection: %v", err)
	}

	res, err := c.(*Conn).Do("CAPABILITIES")

	if err != nil {
		return nil, fmt.Errorf("error getting capabilities: %v", err)
	}

	b := &bytes.Buffer{}

	io.Copy(b, res.Body)

	bs := bytes.TrimSpace(b.Bytes())
	caps := strings.Split(string(bs), "\r\n")

	return caps, nil
}
