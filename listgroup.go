package nntp

import (
	"fmt"
	"io"
)

func (cli *Client) ListGroup(gid string) ([]string, error) {
	res, err := cli.Do("LISTGROUP %s", gid)

	if err != nil {
		return nil, err
	}

	if res.Body == nil {
		return nil, fmt.Errorf("listgroup: no body")
	}
	defer res.Body.Close()

	r := make([]string, 0, 100)

	for {
		var i string

		_, err := fmt.Fscanf(res.Body, "%s\r\n", &i)

		if i != "" {
			r = append(r, i)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
	}

	return r, nil
}
