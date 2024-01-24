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
	TLS                bool

	p *pool.Pool
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

func (c *Client) newConn() (interface{}, error) {
	var conn net.Conn
	var err error
	if c.TLS {
		conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%d", c.Server, c.Port), nil)
	} else {
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", c.Server, c.Port))
	}

	if err != nil {
		return nil, err
	}

	_, nConn, err := NewConn(conn)
	if err != nil {
		return nil, fmt.Errorf("error creating new connection: %v", err)
	}

	if c.User != "" {
		err = nConn.Auth(c.User, c.Pass)

		if err != nil {
			return nil, err
		}
	}

	return nConn, err
}

func NewClient(server string, port int) *Client {
	cli := Client{
		Server: server,
		Port:   port,
	}
	cli.p = pool.NewPool(cli.newConn)
	return &cli
}

func (cli *Client) SetMaxConns(n int) {
	cli.MaxConns = n
	cli.p.SetMax(uint(n))
}
