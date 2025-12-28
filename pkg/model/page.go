package model

import (
	"github.com/oklog/ulid/v2"
)

type Page struct {
	*Search    `json:"-"`
	UserID     ulid.ULID                      `json:"user_id"`
	SearchID   ulid.ULID                      `json:"search_id"`
	ResultID   ulid.ULID                      `json:"result_id"`
	ID         ulid.ULID                      `json:"id"`
	URL        string                         `json:"url"`
	Title      string                         `json:"title"`
	Content    string                         `json:"content"`
	Screenshot string                         `json:"screenshot"`
	Results    map[SearchResultType][]*Result `json:"results"`
}

func (p *Page) Path() string {
	return p.Search.Dir() + "/result"
}

func (p *Page) Dir() string {
	return p.Path() + "/" + p.ID.String()
}

func (p *Page) Key() string {
	return p.Dir() + "/_.json"
}
