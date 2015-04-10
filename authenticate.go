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

	conn := cli.p.Get().(*Conn)
	defer cli.p.Put(conn)

	res, err := conn.Do("AUTHINFO USER %s", u)

	if err != nil {
		return fmt.Errorf("Error authenticating user: %v", err)
	}

	switch res.Code {
	case AuthAccepted:
		return nil
	case PasswordNeeded:
		res, err = conn.Do("AUTHINFO PASSWORD %s", p)

		if res.Code == AuthAccepted {
			err = nil
		}
	}

	return err
}
