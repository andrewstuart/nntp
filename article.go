package nntp

import (
	"io"
	"net/http"
)

const (
	ArticleFound    = 220
	NoArticleWithId = 430
)

type Article struct {
	Headers http.Header
	Body    io.ReadCloser
}

func (a *Article) getHeaders() error {
	//TODO
	return nil
}

func NewArticle(r io.Reader) (*Article, error) {
	//TODO
	return nil, nil
}

func GetArticle(id string) (*Article, error) {
	//TODO
	return nil, nil
}
