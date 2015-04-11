package nntp

import "fmt"

const (
	ArticleFound    = 220
	NoArticleWithId = 430
)

//Client method GetArticle
func (cli *Client) GetArticle(group, id string) (res *Response, err error) {
	conn := cli.p.Get().(*Conn)

	res, err = conn.Do("GROUP %s", group)

	if err != nil {
		return nil, err
	}

	if res.Code != GroupJoined {
		return nil, fmt.Errorf("bad group: %s", res.Message)
	}

	res, err = conn.Do("ARTICLE <%s>", id)

	if err != nil {
		defer cli.p.Put(conn)
		return nil, err
	}

	if res.Code == NoArticleWithId {
		return nil, fmt.Errorf("no article with id %s", id)
	}

	if res.Body != nil {
		//Wraps body in a Closer that returns the connection to the pool.
		res.Body = getPoolBody(cli.p, conn, res.Body)
	} else {
		defer cli.p.Put(conn)
	}

	return res, err
}
