package model

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/s3"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

const ItemPath = "item"

type Items []Item
type Item struct {
	Job    Job       `json:"-" xml:"-"`
	JobID  ulid.ULID `json:"job_id" xml:"-"`
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

func (i Item) IsRedirect() bool {
	return strings.HasPrefix("https://news.google.com/", i.URL) ||
		strings.HasPrefix("https://bing.com/news/", i.URL)
}

func (i Item) Create(db s3.Client, j Job) error {
	i.Job = j
	i.JobID = j.ID
	if err := db.Put(i.Key(), app.MustMarshal(i)); err != nil {
		return err
	}
	return MakeURL(i.URL).Create(db)
}

func (i Item) Path() string {
	return fmt.Sprintf("%s/%s", i.Job.Key(), ItemPath)
}

func (i Item) Key() string {
	return fmt.Sprintf("%s/%s", i.Path(), base64.URLEncoding.EncodeToString([]byte(i.URL)))
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
