package nntp

import "fmt"

const (
	GroupJoined = 211
	NoSuchGroup = 411
)

func (cli *Client) JoinGroup(id string) error {
	r, err := cli.Do("GROUP %s", id)

	if err != nil {
		return fmt.Errorf("error doing group cmd: %v", err)
	}

	switch r.Code {
	case NoSuchGroup:
		return fmt.Errorf("No such group %s. Server says: %s", id, r.Message)
	case GroupJoined:
		cli.CurrGroup = id
		return nil
	}

	return fmt.Errorf("Unexpected response: %s", r)
}
