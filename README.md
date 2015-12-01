# nntp
--
    import "github.com/andrewstuart/nntp"


## Usage

```go
const (
	ArticleFound    = 220
	NoArticleWithId = 430
)
```

```go
const (
	AuthAccepted   = 281
	PasswordNeeded = 381
	AuthNeeded     = 480
	BadAuth        = 481
	ConnsExceeded  = 502
)
```
https://tools.ietf.org/html/rfc4643

```go
const (
	GroupJoined = 211
	NoSuchGroup = 411
)
```

```go
const (
	CapabilitiesFollow = 101
)
```

```go
const HeadersFollow = 221
```

```go
const (
	InfoFollows = 215
)
```

```go
var (
	TooManyConns = ConnErr{ConnsExceeded, "too many connections"}
	AuthRejected = ConnErr{BadAuth, "credentials rejected"}
)
```

```go
var (
	IllegalResponse = fmt.Errorf("illegal response")
	IllegalHeader   = fmt.Errorf("illegal headers")
)
```

#### type Client

```go
type Client struct {
	MaxConns, Port     int
	Server, User, Pass string
	Tls                bool
}
```


#### func  NewClient

```go
func NewClient(server string, port int) *Client
```

#### func (*Client) Auth

```go
func (cli *Client) Auth(u, p string) error
```

#### func (*Client) Capabilities

```go
func (cli *Client) Capabilities() ([]string, error)
```

#### func (*Client) Do

```go
func (cli *Client) Do(format string, args ...interface{}) (*Response, error)
```

#### func (*Client) GetArticle

```go
func (cli *Client) GetArticle(group, id string) (res *Response, err error)
```
Client method GetArticle

#### func (*Client) Head

```go
func (cli *Client) Head(group, id string) (*Response, error)
```

#### func (*Client) JoinGroup

```go
func (cli *Client) JoinGroup(name string) error
```

#### func (*Client) List

```go
func (cli *Client) List() ([]Group, error)
```

#### func (*Client) ListGroup

```go
func (cli *Client) ListGroup(gid string) ([]string, error)
```

#### func (*Client) SetMaxConns

```go
func (cli *Client) SetMaxConns(n int)
```

#### type Conn

```go
type Conn struct {
}
```


#### func  NewConn

```go
func NewConn(c io.ReadWriteCloser, wrappers ...func(io.Reader) io.Reader) *Conn
```

#### func (*Conn) Auth

```go
func (conn *Conn) Auth(u, p string) error
```

#### func (*Conn) Close

```go
func (c *Conn) Close() error
```

#### func (*Conn) Do

```go
func (c *Conn) Do(format string, is ...interface{}) (*Response, error)
```

#### func (*Conn) Wrap

```go
func (c *Conn) Wrap(fn ...func(io.Reader) io.Reader)
```

#### type ConnErr

```go
type ConnErr struct {
	Code   int    `json:"code"xml:"code"`
	Reason string `json:"reason"xml:"reason"`
}
```


#### func (ConnErr) Error

```go
func (c ConnErr) Error() string
```

#### type Group

```go
type Group struct {
	Id           string
	Count, First int
}
```


#### type Reader

```go
type Reader struct {
	R *bufio.Reader
}
```

A Reader is a read/closer that strips NNTP newlines and will unescape
characters.

#### func  NewReader

```go
func NewReader(r io.Reader) *Reader
```
NewReader returns an nntp.Reader for the body of the nttp article.

#### func (*Reader) Close

```go
func (r *Reader) Close() error
```
Close allows users of a Reader to signal that they are done using the reader.

#### func (*Reader) Read

```go
func (r *Reader) Read(p []byte) (bytesRead int, err error)
```
The Read method handles translation of the NNTP escaping and marking EOF when
the end of a body is received.

#### type Response

```go
type Response struct {
	Code    int                  `json:"code"xml:"code"`
	Message string               `json:"message"xml:"message"`
	Headers textproto.MIMEHeader `json:"headers"xml:"headers"`
	Body    io.ReadCloser        `json:"body"xml:"body"` //Presence (non-nil) indicates multi-line response
}
```


#### func  NewResponse

```go
func NewResponse(r io.Reader) (*Response, error)
```
