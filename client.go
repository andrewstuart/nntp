package nntp

import (
	"fmt"
	"io"
	"net"

	"git.astuart.co/andrew/pool"
)

type Client struct {
	MaxConns, Port     int
	Server, User, Pass string

	nConns int
	p      *pool.Pool

	cls chan (chan error)
}

func (cli *Client) Do(format string, args ...interface{}) (*Response, error) {
	conn := cli.p.Get().(*Conn)

	res, err := conn.Do(format, args...)

	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		res.Body = getPoolBody(cli.p, conn, res.Body)
	} else {
		cli.p.Put(conn)
	}

	return res, nil
}

func getPoolBody(p pool.Pooler, conn *Conn, rc io.ReadCloser) *poolBody {
	return &poolBody{
		ReadCloser: rc,
		p:          p,
		conn:       conn,
	}
}

type poolBody struct {
	io.ReadCloser
	p    pool.Pooler
	conn *Conn
}

func (pb *poolBody) Close() error {
	pb.p.Put(pb.conn)
	return pb.ReadCloser.Close()
}

func NewClient(server string, port, conns int) *Client {
	cli := Client{
		Server:   server,
		Port:     port,
		MaxConns: conns,
	}

	makeConn := pool.NewFunc(func() (interface{}, error) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cli.Server, cli.Port))

		if err != nil {
			conn.Close()
			return nil, err
		}

		nConn := NewConn(conn)

		if cli.User != "" {
			err = nConn.Auth(cli.User, cli.Pass)

			if err != nil {
				return nil, err
			}
		}

		return nConn, err
	})

	cli.p = pool.NewPool(makeConn)
	cli.p.SetMax(uint(conns))

	return &cli
}
