package model

import (
	"encoding/xml"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type RSS struct {
	Channel struct {
		Items Items `xml:"item"`
	} `xml:"channel"`
}

type Items []Item

func (ii Items) MarshalZerologArray(a *zerolog.Array) {
	for _, i := range ii {
		a.Object(i)
	}
}

type Item struct {
	Title  string    `json:"title" xml:"title"`
	Link   string    `json:"link" xml:"link"`
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

func (i Item) MarshalZerologObject(e *zerolog.Event) {
	e.Str("title", i.Title).Str("url", i.Link)
}

type DateTime time.Time

func (v *DateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	s = strings.Trim(s, `"`) // Remove quotes from the JSON string
	if s == "" || s == "null" {
		return nil // Handle empty or null strings
	}

	t, err := time.Parse(time.RFC1123, s) // Parse using your custom format
	if err != nil {
		return err
	}

	*v = DateTime(t)
	return nil
}

func (v *DateTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = strings.Trim(s, `"`) // Remove quotes from the JSON string
	if s == "" || s == "null" {
		return nil // Handle empty or null strings
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	*v = DateTime(t)
	return nil
}

func (v *DateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(*v).Format(time.RFC3339) + `"`), nil
}
