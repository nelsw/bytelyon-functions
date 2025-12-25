package internal

import (
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
)

type Page struct {
	UserID   ulid.ULID `json:"user_id"`
	SearchID ulid.ULID `json:"search_id"`
	ResultID ulid.ULID `json:"result_id"`
	ID       ulid.ULID `json:"id"`
	URL      string    `json:"url"`
	Title    string    `json:"title"`
}

func (p *Page) MarshalZerologObject(evt *zerolog.Event) {
	evt.Stringer("page", p).
		Str("title", p.Title)
}

func (p *Page) String() string {
	return Path("user", p.UserID, "search", p.SearchID, "result", p.ResultID, "page", p.ID)
}
