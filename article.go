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
	conn.grp = group

	if err != nil {
		return nil, err
	}

	if res.Code != GroupJoined {
		return nil, fmt.Errorf("bad group: %s", res.Message)
	}

	res, err = conn.Do("ARTICLE %s", id)

	if err != nil {
		defer cli.p.Put(conn)
		return nil, err
	}

	if res.Body != nil {
		res.Body = getPoolBody(cli.p, conn, res.Body)
	} else {
		defer cli.p.Put(conn)
	}

	return res, err
}
