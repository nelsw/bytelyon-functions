package model

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var bingRegexp *regexp.Regexp

func init() {
	bingRegexp = regexp.MustCompile("</?News(:\\w+)>")
}

type Work struct {
	Job   Job       `json:"-"`
	ID    ulid.ULID `json:"id"`
	Items Items     `json:"items"`
	Err   error     `json:"error"`
}

func (w Work) Key() string {
	return fmt.Sprintf("%s/work/%s.json", w.Job.Path(), w.ID)
}

func NewWork(j Job) Work {

	w := Work{ID: NewUlid()}

	switch j.Type {
	case NewsJobType:
		w.Items, w.Err = workNews(j)
	default:
		w.Err = errors.New(fmt.Sprintf("unknown job type [%d]", j.Type))
	}

	return w
}

func workNews(j Job) (items Items, err error) {

	f := func(u string) (Items, error) {
		u = fmt.Sprintf(u, strings.Join(j.Keywords, ","))
		res, e := http.Get(u)
		if e != nil {
			return nil, errors.Join(errors.New("failed to http.Get url"), e)
		}
		defer res.Body.Close()

		var b []byte
		if b, e = io.ReadAll(res.Body); e != nil {
			return nil, errors.Join(errors.New("failed to io.ReadAll response body"), e)
		}

		if strings.Contains(u, "https://www.bing.com/") {
			str := bingRegexp.ReplaceAllStringFunc(string(b), func(s string) string {
				return strings.ReplaceAll(s, ":", "_")
			})
			b = []byte(str)
		}

		var rss RSS
		if e = xml.Unmarshal(b, &rss); e != nil {
			return nil, errors.Join(errors.New("failed to unmarshal xml"), e)
		}

		return rss.Channel.Items, nil
	}

	for _, u := range j.Type.URLs() {
		ii, e := f(u)
		err = errors.Join(err, e)
		if ii != nil {
			items = append(items, ii...)
		}
		log.Err(e).Any("job", j).Array("items", ii).Send()
	}

	return
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

func (v *DateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(*v).Format(time.RFC3339) + `"`), nil
}

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
