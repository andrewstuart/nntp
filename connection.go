package nntp

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

//a connection is a buffered reader and unbuffered writer
type connection struct {
	io.Writer
	br *bufio.Reader
}

func newConnection(rw io.ReadWriter) *connection {
	return &connection{
		io.Writer(rw),
		bufio.NewReader(rw),
	}
}

//init sets up a connection to the server
func (cli *Client) getConn() (*connection, error) {
	select {
	case c := <-cli.cBucket:
		return c, nil
	default:
		if len(cli.conns) >= cli.Connections {
			return <-cli.cBucket, nil
		}
	}

	//If a connection wasn't already available and we aren't yet over our limit,
	//make a new connection and return it

	server := fmt.Sprintf("%s:%d", cli.Server, cli.Port)
	conn, err := net.Dial("tcp", server)

	if err != nil {
		return nil, fmt.Errorf("tcp error: %v", err)
	}

	bufCon := newConnection(conn)

	//Drop hello
	_, err = bufCon.br.ReadBytes('\n')

	if err != nil {
		return nil, fmt.Errorf("error reading WELCOME message: %v", err)
	}

	//Automatically authenticate new connections
	err = bufCon.Auth(cli.Username, cli.Password)

	if err != nil {
		if err == TooManyConnections {
			return <-cli.cBucket, nil
		}

		return nil, fmt.Errorf("error authenticating: %v", err)
	}

	if cli.CurrGroup != "" {
		res, err := bufCon.do("GROUP %s", cli.CurrGroup)

		if err != nil {
			return nil, fmt.Errorf("error connecting to group: %v", err)
		}

		if res.Code != GroupJoined {
			return nil, fmt.Errorf("could not join group %s: %v", cli.CurrGroup, res.Message)
		}
	}

	cli.conns = append(cli.conns, bufCon)
	return bufCon, nil
}
