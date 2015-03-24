package nntp

import "fmt"

const (
	//https://tools.ietf.org/html/rfc4643
	AuthAccepted   = 281
	PasswordNeeded = 381
	AuthRejected   = 481
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
		case AuthAccepted:
			return nil
		case AuthRejected:
			return fmt.Errorf("Authentication Rejected")
		}
	}

	return nil
}
