package internal

import "github.com/rs/zerolog"

type Result struct {
	Position int     `json:"position"`
	Title    string  `json:"title"`
	Link     string  `json:"link"`
	Source   string  `json:"source"`
	Snippet  string  `json:"snippet"`
	Price    float64 `json:"price,omitempty"`
}

func (r *Result) MarshalZerologObject(evt *zerolog.Event) {
	evt.Int("pos", r.Position).
		Str("title", r.Title).
		Str("link", r.Link).
		Str("src", r.Source)

	if r.Snippet != "" {
		evt.Str("snippet", r.Snippet)
	}

	if r.Price > 0 {
		evt.Float64("price", r.Price)
	}
}
