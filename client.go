package nntp

//A Client is the key abstraction for connecting to a server
type Client struct {
	Server, Username, Password, CurrGroup string
	Port, Connections, Retention, Timeout int
	cBucket                               chan *connection
	conns                                 []*connection
	nConns                                int
	compression                           bool
}

//NewClient returns a pointer to a downloader
func NewClient(s string, port, conns int) *Client {
	if conns == 0 {
		conns = 1
	}

	cli := Client{
		Server:      s,
		Port:        port,
		Connections: conns,
		cBucket:     make(chan *connection, conns),
		conns:       make([]*connection, 0, conns),
	}

	return &cli
}
