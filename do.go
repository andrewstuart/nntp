package nntp

import (
	"fmt"
	"strconv"
	"strings"
)

type Response struct {
	Code    int
	Message string
}

func (r *Response) String() string {
	return fmt.Sprintf("%d: %s", r.Code, r.Message)
}

func (c *connection) do(cmd string, args ...interface{}) (*Response, error) {
	_, err := fmt.Fprintf(c, cmd+"\r\n", args...)

	if err != nil {
		return nil, err
	}

	s, err := c.br.ReadString('\n')
	if err != nil {
		return nil, err
	}

	ss := strings.SplitN(strings.TrimSpace(s), " ", 2)

	r := Response{}

	if len(ss) > 1 {
		r.Message = ss[1]
		r.Code, err = strconv.Atoi(ss[0])

		if err != nil {
			return nil, fmt.Errorf("error casting error code: %v", err)
		}
	}

	return &r, nil
}

func (cli *Client) Do(cmd string, args ...interface{}) (*Response, error) {
	//Get a connection from the pool
	conn, err := cli.getConn()
	//Don't forget to put it back
	defer func() {
		cli.cBucket <- conn
	}()

	if err != nil {
		return nil, fmt.Errorf("error making connection: %v", err)
	}

	//Do your stuff
	res, err := conn.do(cmd, args...)

	if err != nil {
		return nil, fmt.Errorf("error executing command: %v", err)
	}

	return res, nil
}
