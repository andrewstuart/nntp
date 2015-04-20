package nntp

import "fmt"

const HeadersFollow = 221

func (cli *Client) Head(group, id string) (*Response, error) {
	if err := cli.JoinGroup(group); err != nil {
		return nil, fmt.Errorf("Could not join group %s: %v", group, err)
	}

	res, err := cli.Head("a", "123")

	if res.Body == nil {
		return nil, fmt.Errorf("no header body")
	}

	if res.Code != HeadersFollow {
		return nil, fmt.Errorf("error getting headers: %s", res.Message)
	}

	defer res.Body.Close()

	return res, err
}
