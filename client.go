package nntp

import "sync"

type Client struct {
	MaxConns, Port     int
	Server, User, Pass string

	p *sync.Pool
}
