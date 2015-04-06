package nntp

type Client struct {
	MaxConns, Port     int
	Server, User, Pass string
}
