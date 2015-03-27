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
	case c := <-cli.pool:
		return c, nil
	default:
		if cli.nConns >= cli.Connections {
			return <-cli.pool, nil
		}
	}

	cli.nConns++

	//If a connection wasn't already available and we aren't yet over our limit,
	//make a new connection and return it

	server := fmt.Sprintf("%s:%d", cli.Server, cli.Port)
	conn, err := net.Dial("tcp", server)

	if err != nil {
		cli.nConns--
		return nil, fmt.Errorf("tcp error: %v", err)
	}

	bufCon := newConnection(conn)

	//Drop hello
	_, err = bufCon.br.ReadBytes('\n')

	if err != nil {
		cli.nConns--
		return nil, fmt.Errorf("error reading WELCOME message: %v", err)
	}

	//Automatically authenticate new connections
	err = bufCon.Auth(cli.Username, cli.Password)

	if err != nil {
		cli.nConns--
		if err == TooManyConnections {
			return <-cli.pool, nil
		}

		return nil, fmt.Errorf("error authenticating: %v", err)
	}

	if cli.CurrGroup != "" {
		res, err := bufCon.do("GROUP %s", cli.CurrGroup)

		if err != nil {
			cli.nConns--
			return nil, fmt.Errorf("error connecting to group: %v", err)
		}

		if res.Code != GroupJoined {
			cli.nConns--
			return nil, fmt.Errorf("could not join group %s: %v", cli.CurrGroup, res.Message)
		}
	}

	return bufCon, nil
}
