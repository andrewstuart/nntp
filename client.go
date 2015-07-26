package nntp

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"time"

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

//The unexported newConn function is used by the client's connection pool to
//create a new wrapped tcp connection when possible.
func (c *Client) newConn() (interface{}, error) {
	var conn io.ReadWriteCloser
	var err error
	if c.Tls {
		conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%d", c.Server, c.Port), nil)
	} else {
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", c.Server, c.Port))
	}

	if err != nil {
		return nil, err
	}

	nConn := NewConn(conn)

	if c.User != "" {
		err = nConn.Auth(c.User, c.Pass)

		if err != nil {
			return nil, err
		}
	}

	return nConn, err
}

func (c *Client) SetTimeout(d time.Duration) {
	c.p.SetTimeout(d)
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
