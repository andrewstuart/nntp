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

	connChan chan (chan *Conn)
	nConns   int
	p        *pool.Pool

	cls chan (chan error)
}

func (cli *Client) Do(format string, args ...interface{}) (*Response, error) {
	conn := cli.p.Get().(*Conn)

	res, err := conn.Do(format, args...)

	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		res.Body = &poolBody{
			ReadCloser: res.Body,
			cli:        cli,
			conn:       conn,
		}
	} else {
		cli.p.Put(conn)
	}

	return res, nil
}

type poolBody struct {
	io.ReadCloser
	cli  *Client
	conn *Conn
}

func (pb *poolBody) Close() error {
	pb.cli.p.Put(pb.conn)
	return pb.ReadCloser.Close()
}

func NewClient(server string, port, conns int) *Client {
	cli := Client{
		Server:   server,
		Port:     port,
		MaxConns: conns,
		connChan: make(chan (chan *Conn)),
	}

	makeConn := pool.NewFunc(func() (interface{}, error) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cli.Server, cli.Port))

		if err != nil {
			return nil, err
		}

		return NewConn(conn), nil
	})

	cli.p = pool.NewPool(makeConn)

	return &cli
}
