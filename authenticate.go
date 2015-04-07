package nntp

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

	res, err := conn.Do("AUTHINFO USER %s", u)

	switch res.Code {
	case AuthAccepted:
		return nil
	case PasswordNeeded:
		res, err = conn.Do("AUTHINFO PASSWORD %s", p)
	default:
	}

	return err
}
