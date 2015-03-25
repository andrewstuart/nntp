package nntp

import (
	"fmt"
	"io/ioutil"
)

const (
	InfoFollows = 215
)

func (cli *Client) List(search string) ([]byte, error) {
	conn, err := cli.getConn()

	if err != nil {
		return nil, fmt.Errorf("could not get connection")
	}

	res, err := conn.do("LIST %s", search)

	switch res.Code {
	case InfoFollows:
		return ioutil.ReadAll(conn.br)
	}

	return nil, fmt.Errorf("unexpected message")
}
