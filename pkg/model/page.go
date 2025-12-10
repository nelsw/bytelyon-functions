package model

import (
	"bytelyon-functions/pkg/service/s3"
	"encoding/json"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Page struct {
	*Search    `json:"-"`
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
	evt.Stringer("id", p.ID).
		Str("url", p.URL).
		Str("title", p.Title).
		Bool("content", p.Content != "").
		Bool("screenshot", p.Screenshot != "").
		Bool("results", p.Results != nil)
}

func (p *Page) Path() string {
	return p.Search.Dir() + "/page"
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
