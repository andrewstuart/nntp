package nntp

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"

	"github.com/andrewstuart/pool"
)

type Client struct {
	MaxConns, Port     int
	Server, User, Pass string
	Tls                bool

	nConns int
	p      *pool.Pool

	cls chan (chan error)
}

func (cli *Client) Do(format string, args ...interface{}) (*Response, error) {
	c, err := cli.p.Get()

	if err != nil {
		return nil, fmt.Errorf("error getting client: %v", err)
	}

	conn := c.(*Conn)

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

func NewClient(server string, port int) *Client {
	cli := Client{
		Server: server,
		Port:   port,
	}

	cpt := &cli

	makeConn := pool.NewFunc(func() (interface{}, error) {
		var conn io.ReadWriteCloser
		var err error
		if cpt.Tls {
			conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%d", cpt.Server, cpt.Port), nil)
		} else {
			conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", cpt.Server, cpt.Port))
		}

		if err != nil {
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

	return cpt
}

func (cli *Client) SetMaxConns(n int) {
	cli.MaxConns = n
	cli.p.SetMax(uint(n))
}
