package nntp

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	ArticleFound    = 220
	NoArticleWithId = 430
)

type Article struct {
	Headers map[string]string
	Body    io.Reader
}

func (a *Article) getHeaders() error {
	b := bufio.NewReader(a.Body)

	if a.Headers == nil {
		a.Headers = make(map[string]string)
	}

	for {
		s, err := b.ReadString('\n')

		if err != nil {
			return err
		}

		s = strings.TrimSpace(s)

		//Blank line inidicates start of body
		if s == "" {
			return nil
		}

		hparts := strings.Split(s, ": ")

		if len(hparts) > 1 {
			a.Headers[hparts[0]] = hparts[1]
		}
	}
	a.Body = b

	return nil
}

func NewArticle(r io.Reader) (*Article, error) {
	a := Article{make(map[string]string), bufio.NewReader(r)}

	if err := a.getHeaders(); err != nil {
		return nil, fmt.Errorf("error reading headers: %v", err)
	}

	return &a, nil
}

//GetArticle returns a reader for the article, or an error indicating the reason an article
//could not be founc.
func (cli *Client) GetArticle(id string) (*Article, error) {
	conn, err := cli.getConn()

	if err != nil {
		return nil, fmt.Errorf("error getting connection from pool: %v", err)
	}

	res, err := conn.do("ARTICLE <%s>", id)
	if err != nil {
		return nil, err
	}

	switch res.Code {
	case ArticleFound:
		bdy := NewArticleReader(conn.br)

		art, err := NewArticle(bdy)

		if err != nil {
			return nil, fmt.Errorf("header error: %v", err)
		}

		if b, isBody := bdy.(*body); isBody {
			go func() {
				b.done.Wait()
				cli.cBucket <- conn
			}()
		}

		return art, nil
	case NoArticleWithId:
		cli.cBucket <- conn
		return nil, fmt.Errorf("no article with ID %s founc. Server says: %v", id, res)
	}

	cli.cBucket <- conn
	return nil, fmt.Errorf("unexpected response: %s", res)
}
