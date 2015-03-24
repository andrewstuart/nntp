package nntp

import "fmt"

type connErr struct {
	code   int
	reason string
}

func (c connErr) Error() string {
	return fmt.Sprintf("%d: %s", c.code, c.reason)
}

const (
	//https://tools.ietf.org/html/rfc4643
	AuthAccepted   = 281
	PasswordNeeded = 381
	AuthNeeded     = 480
	AuthRejected   = 481
)

var (
	TooManyConnections = connErr{502, "Too many connections"}
)

func (c *connection) Auth(user, pass string) error {
	//Check for username
	if user == "" {
		return fmt.Errorf("No username specified")
	}

	r, err := c.do("AUTHINFO USER %s", user)

	if err != nil {
		return err
	}

	if r.Code == PasswordNeeded {
		if pass == "" {
			return fmt.Errorf("Password needed for user %s and was not set.", user)
		}

		r, err := c.do("AUTHINFO PASS %s", pass)

		if err != nil {
			return err
		}

		switch r.Code {
		case TooManyConnections.code:
			return TooManyConnections
		case AuthAccepted:
			return nil
		case AuthRejected:
			return fmt.Errorf("Authentication Rejected")
		default:
			return fmt.Errorf("Unexpected code: %v", r)
		}
	}

	return nil
}
