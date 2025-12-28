package model

import (
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Item struct {
	ID     ulid.ULID `json:"id" xml:"-"`
	URL    string    `json:"link" xml:"link"`
	Title  string    `json:"title" xml:"title"`
	Time   *DateTime `json:"published" xml:"pubDate"`
	Source *struct {
		URL   string `json:"url,omitempty" xml:"url,attr"`
		Value string `json:"value,omitempty" xml:",chardata"`
	} `json:"source,omitempty" xml:"source"`
	NewsSource         string `json:"news_source,omitempty" xml:"News_Source"`
	NewsImage          string `json:"news_image,omitempty" xml:"News_Image"`
	NewsImageSize      string `json:"news_image_size,omitempty" xml:"News_ImageSize"`
	NewsImageMaxWidth  int    `json:"news_image_max_width,omitempty" xml:"News_ImageMaxWidth"`
	NewsImageMaxHeight int    `json:"news_image_max_height,omitempty" xml:"News_ImageMaxHeight"`
}

func (i *Item) MarshalZerologObject(evt *zerolog.Event) {
	evt.Stringer("id", i.ID).
		Str("url", i.URL).
		Str("title", i.Title).
		Time("time", time.Time(*i.Time))
	if i.Source != nil {
		evt.Str("source", i.Source.Value)
	}
	evt.Str("news_source", i.NewsSource)
}

func (i *Item) Scrub() {

	i.ID = i.Time.ULID()

	if strings.HasPrefix(i.URL, "https://news.google.com/") {
		if s, err := decodeGoogleNewsURL(i.URL); err != nil {
			log.Warn().Err(err).Msg("failed to decode google news url")
		} else {
			i.URL = s
		}
		return
	}

	if strings.HasPrefix(i.URL, "http://www.bing.com/") {
		if s, err := decodeBingNewsURL(i.URL); err != nil {
			log.Warn().Err(err).Msg("failed to decode bing news url")
		} else {
			i.URL = s
		}
	}
}
