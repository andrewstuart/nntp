package nntp

import (
	"bytes"
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
	} else if err != nil {
		return nil, err
	}

	if res.Body != nil {
		//Drain body if not nil
		defer res.Body.Close()
		io.Copy(&bytes.Buffer{}, res.Body)
	}

	if res.Code != HeadersFollow {
		return nil, fmt.Errorf("error getting headers: %s", res.Message)
	}

	return res, err
}
