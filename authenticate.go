package nntp

import "fmt"

type ConnErr struct {
	Code   int    `json:"code"xml:"code"`
	Reason string `json:"reason"xml:"reason"`
}

func (c ConnErr) Error() string {
	return fmt.Sprintf("%d: %s", c.Code, c.Reason)
}

const (
	//https://tools.ietf.org/html/rfc4643
	AuthAccepted   = 281
	PasswordNeeded = 381
	AuthNeeded     = 480
	AuthRejected   = 481
	TooManyConns   = 502
)

var (
	ErrTooManyConns = ConnErr{TooManyConns, "too many connections"}
	ErrAuthRejected = ConnErr{AuthRejected, "credentials rejected"}
)
