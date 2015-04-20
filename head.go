package nntp

import (
	"fmt"
	"io"
)

const HeadersFollow = 221

func (cli *Client) Head(group, id string) (*Response, error) {
	if err := cli.JoinGroup(group); err != nil {
		return nil, fmt.Errorf("Could not join group %s: %v", group, err)
	}

	res, err := cli.Do("HEAD %s", id)

	if err == io.EOF && res != nil {
		return res, nil
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.Code != HeadersFollow {
		return nil, fmt.Errorf("error getting headers: %s", res.Message)
	}

	return res, err
}
