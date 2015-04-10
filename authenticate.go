package nntp

import "fmt"

//https://tools.ietf.org/html/rfc4643
const (
	AuthAccepted   = 281
	PasswordNeeded = 381
	AuthNeeded     = 480
	BadAuth        = 481
	ConnsExceeded  = 502
)

var (
	TooManyConns = ConnErr{ConnsExceeded, "too many connections"}
	AuthRejected = ConnErr{BadAuth, "credentials rejected"}
)

func (cli *Client) Auth(u, p string) error {
	cli.User = u
	cli.Pass = p
	return nil
}

func (conn *Conn) Auth(u, p string) error {
	res, err := conn.Do("AUTHINFO USER %s", u)

	if err != nil {
		return fmt.Errorf("error authenticating user: %v", err)
	}

	switch res.Code {
	case AuthAccepted:
		return nil
	case PasswordNeeded:
		res, err = conn.Do("AUTHINFO PASS %s", p)
	}

	return err
}
