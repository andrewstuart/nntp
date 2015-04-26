package nntp

import (
	"fmt"
	"io"
)

const (
	InfoFollows = 215
)

func (cli *Client) List() ([]Group, error) {
	res, err := cli.Do("LIST ACTIVE")

	if err != nil {
		return nil, err
	}

	groups := make([]Group, 0)

	s := ""
	for {
		g := Group{}
		_, err := fmt.Fscanf(res.Body, "%s %d %d %s\r\n", &g.Id, &g.Count, &g.First, &s)

		if g.Id != "" {
			groups = append(groups, g)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}

	return groups, nil
}
