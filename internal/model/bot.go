package model

import (
	"bytelyon-functions/internal/entity"
	"encoding/xml"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

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

type JobType int

const (
	GoogleNews JobType = iota
	BingNews
)

func (v JobType) URLs() []string {
	switch v {
	case GoogleNews:
		return []string{"https://news.google.com/rss/search?q=%s&hl=en-US&gl=US&ceid=US:en"}
	case BingNews:
		return []string{"https://www.bing.com/search?format=rss&q=", "https://www.bing.com/news/search?format=rss&q="}
	default:
		return nil
	}
}

func (v JobType) Sanitize(b []byte) []byte {

	if v != BingNews {
		return b
	}

	reg := regexp.MustCompile("</?News(:\\w+)>")
	str := reg.ReplaceAllStringFunc(string(b), func(s string) string {
		return strings.ReplaceAll(s, ":", "_")
	})

	return []byte(str)
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

type RSS struct {
	Channel struct {
		Items []Item `xml:"item"`
	} `xml:"channel"`
}
type Work struct {
	ID    ulid.ULID `json:"id"`
	JobID uuid.UUID `json:"job_id"`
	Items []Item    `json:"items"`
}

type Job struct {
	ID        uuid.UUID `json:"id"`
	Type      JobType   `json:"type"`
	Keywords  []string  `json:"keywords"`
	Frequency int       `json:"frequency"` // minutes
	LastRun   time.Time `json:"last_run"`
}

func (v *Job) ReadyToWork() bool {
	return v.LastRun.Add(time.Minute * time.Duration(v.Frequency)).Before(time.Now())
}

func (v *Job) CreateWork() error {

	fn := func(u string) {

		res, err := http.Get(u)
		if err != nil {
			log.Error().Err(err).Str("URL", u).Str("ID", v.ID.String()).Msg("failed to http.Get url")
			return
		}
		defer res.Body.Close()

		var b []byte
		if b, err = io.ReadAll(res.Body); err != nil {
			log.Error().Err(err).Str("ID", v.ID.String()).Msg("failed to io.ReadAll response body")
			return
		}

		b = v.Type.Sanitize(b)

		var rss RSS
		if err = xml.Unmarshal(b, &rss); err != nil {
			log.Error().Err(err).Str("ID", v.ID.String()).Msg("failed to unmarshal xml")
			return
		}

		log.Info().Int("size", len(rss.Channel.Items)).Msg("work items found")

		if len(rss.Channel.Items) == 0 {
			return
		}

		work := Work{NewULID(), v.ID, rss.Channel.Items}
		if err = entity.New().Value(&work).Save(); err != nil {
			return
		}

		log.Info().Str("workID", work.ID.String()).Msg("work items created")
	}

	log.Info().Str("ID", v.ID.String()).Msg("creating work")

	for _, u := range v.Type.URLs() {
		fn(u)
	}

	if err := entity.New().Value(v).Save(); err != nil {
		log.Error().Err(err).Str("ID", v.ID.String()).Msg("failed to update job last run time")
		return err
	}

	log.Info().Str("ID", v.ID.String()).Msg("job updated")
	return nil
}
