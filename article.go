package nntp

const (
	ArticleFound    = 220
	NoArticleWithId = 430
)

//Client method GetArticle
func (cli *Client) GetArticle(id string) (*Response, error) {
	return cli.Do("ARTICLE %s", id)
}
