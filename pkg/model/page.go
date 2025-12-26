package model

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
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

func NewPage(a ...any) *Page {
	p := &Page{
		UserID:     a[0].(ulid.ULID),
		SearchID:   a[1].(ulid.ULID),
		ResultID:   a[2].(ulid.ULID),
		ID:         NewUlid(),
		URL:        a[3].(string),
		Title:      a[4].(string),
		Content:    a[5].(string),
		Screenshot: a[6].(string),
	}
	p.HandleResults()
	return p
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

func (p *Page) HandleResults() {
	if !strings.HasPrefix(p.URL, "https://www.google.com") {
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(p.Content))
	if err != nil {
		log.Warn().Err(err).Msg("Page - failed to parse html")
		return
	}
}
