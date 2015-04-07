package nntp

import (
	"fmt"
	"net"
	"sync"
)

type Client struct {
	MaxConns, Port     int
	Server, User, Pass string

	connChan chan (chan *Conn)
	nConns   int
	p        *sync.Pool

	cls chan (chan error)
}

func (cli *Client) run() {
	for {
		select {
		case nc := <-cli.connChan:
			if cli.nConns < cli.MaxConns {
				cli.nConns++
				conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cli.Server, cli.Port))

				if err != nil {
					cli.nConns--
					nc <- nil
				}
				nc <- NewConn(conn)
			} else {
				nc <- nil
			}
		}
	}
}

func (cli *Client) Do(format string, args ...interface{}) (*Response, error) {
	conn := cli.p.Get().(*Conn)

	if conn == nil {
		return nil, fmt.Errorf("client err")
	}

	res, err := conn.Do(format, args...)

	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		conn.onClose = func() error {
			cli.p.Put(conn)
			return nil
		}
	} else {
		cli.p.Put(conn)
	}

	return res, nil
}

func NewClient(server string, port, conns int) *Client {
	cli := Client{
		Server:   server,
		Port:     port,
		MaxConns: conns,
		connChan: make(chan (chan *Conn)),
	}

	go cli.run()

	cli.p = &sync.Pool{
		New: func() interface{} {
			nch := make(chan *Conn)
			cli.connChan <- nch

			conn := <-nch
			if conn == nil {
				cli.p.New = nil
			}

			return conn
		},
	}

	return &cli
}
