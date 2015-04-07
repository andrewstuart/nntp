package nntp

import "fmt"

const (
	GroupJoined = 211
	NoSuchGroup = 411
)

func (cli *Client) JoinGroup(name string) error {
	req, err := cli.Do("GROUP %s", name)

	if err != nil {
		return err
	}

	switch req.Code {
	case GroupJoined:
		return nil
	case NoSuchGroup:
		return fmt.Errorf("no such group: %s", req.Message)
	default:
		return fmt.Errorf(req.Message)
	}
}
