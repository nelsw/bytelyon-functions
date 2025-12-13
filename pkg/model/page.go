package model

import (
	"bytelyon-functions/pkg/service/s3"
	"encoding/json"
	"strings"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Page struct {
	*Search    `json:"-"`
	UserID     ulid.ULID `json:"user_id"`
	SearchID   ulid.ULID `json:"search_id"`
	ResultID   ulid.ULID `json:"result_id"`
	ID         ulid.ULID `json:"id"`
	URL        string    `json:"url"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Screenshot string    `json:"screenshot"`
	Results    *Results  `json:"results"`
}

func NewPage(s *Search) *Page {
	return &Page{Search: s}
}

func (p *Page) MarshalZerologObject(evt *zerolog.Event) {
	if p.Search != nil {
		evt.EmbedObject(p.Search)
	}
	url := strings.TrimPrefix(p.URL, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "www.")
	idx := strings.Index(url, "/")
	if idx > 0 {
		url = url[:idx]
	}
	if len(url) > 8 {
		url = url[:8]
	}
	evt.Str("result", p.ResultID.String()[20:]).
		Str("page", p.ID.String()[20:]).
		Str("search", p.SearchID.String()[20:]).
		Str("user", p.UserID.String()[20:]).
		Str("domain", url).
		Bool("content", p.Content != "").
		Bool("screenshot", p.Screenshot != "").
		Bool("results", p.Results != nil)
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

func (p *Page) FetchData(s *Search) *Page {
	p.Search = s
	db := s3.New()

	if out, err := db.GetPresigned(p.Dir() + "/content.html"); err == nil {
		p.Content = out
	}

	if out, err := db.GetPresigned(p.Dir() + "/screenshot.png"); err == nil {
		p.Screenshot = out
	}

	if out, err := db.Get(p.Dir() + "/results.json"); err == nil {
		if err = json.Unmarshal(out, &p.Results); err != nil {
			log.Warn().Err(err).Msg("failed to unmarshal pagee results")
		}
	}

	return p
}
