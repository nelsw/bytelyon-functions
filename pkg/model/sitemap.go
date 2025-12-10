package model

import (
	"bytelyon-functions/pkg/service/em"
	"bytelyon-functions/pkg/service/s3"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Sitemap struct {
	*User    `json:"-"`
	ID       ulid.ULID `json:"id"`
	URL      string    `json:"url"`
	Domain   string    `json:"domain"`
	Duration float64   `json:"duration"`
	Relative []string  `json:"relative"`
	Remote   []string  `json:"remote"`
}

func (s *Sitemap) MarshalZerologObject(evt *zerolog.Event) {
	if s.User != nil {
		evt.EmbedObject(s.User)
	}
	evt.Stringer("sitemap", s.ID).
		Str("url", s.URL).
		Str("domain", s.Domain).
		Float64("duration", s.Duration).
		Int("relative", len(s.Relative)).
		Int("remote", len(s.Remote))
}

func (s *Sitemap) Path() string {
	return s.User.Dir() + "/sitemap"
}

func (s *Sitemap) Key() string {
	return s.Path() + "/" + s.Domain + "/" + s.ID.String() + "/_.json"
}

func NewSitemap(user *User) *Sitemap {
	return &Sitemap{User: user}
}

func (s *Sitemap) Fetch(url string) ([]string, []string, error) {

	if strings.HasSuffix(url, ".jpeg") ||
		strings.HasSuffix(url, ".png") ||
		strings.HasSuffix(url, ".gif") ||
		strings.HasSuffix(url, ".jpg") ||
		strings.HasSuffix(url, ".pdf") {
		return nil, nil, nil
	}

	doc := NewDocument(url)
	if err := doc.Fetch(); err != nil {
		return nil, nil, err
	}

	var relative, remote []string
	for _, a := range doc.anchors() {

		if regexp.MustCompile(`^(#|mailto:|tel:).*`).MatchString(a) {
			continue
		}

		if strings.HasPrefix(a, "?") {
			relative = append(relative, url+a)
			continue
		}

		if strings.HasPrefix(a, "/") {
			url = strings.TrimSuffix(url, "/")
			relative = append(relative, url+a)
			continue
		}

		a = strings.TrimPrefix(a, "https://")
		a = strings.TrimPrefix(a, "http://")
		a = strings.TrimPrefix(a, "www.")
		a = strings.TrimSpace(a)

		if !strings.HasPrefix(a, s.Domain) {
			remote = append(remote, a)
			continue
		}

		a = strings.TrimPrefix(a, s.Domain)
		relative = append(relative, s.URL+a)
	}

	return relative, remote, nil
}

func (s *Sitemap) Create(b []byte) (*Sitemap, error) {

	log.Info().EmbedObject(s.User).Msg("creating sitemap")

	user := s.User
	if err := json.Unmarshal(b, s); err != nil {
		log.Err(err).Msg("failed to unmarshal sitemap")
		return nil, err
	}

	s.User = user
	s.URL = strings.TrimSpace(s.URL)

	if !strings.HasPrefix(s.URL, "https://") {
		return nil, errors.New("invalid URL")
	}

	if strings.HasSuffix(s.URL, "/") {
		s.URL = strings.TrimSuffix(s.URL, "/")
	}

	s.Domain = strings.TrimPrefix(s.URL, "https://")
	s.Domain = strings.TrimPrefix(s.Domain, "http://")
	s.Domain = strings.TrimPrefix(s.Domain, "www.")

	// new up a Crawler using a reference to the Sitemap, aka Fetcher
	crawler := NewCrawler(s)

	// increment the crawler wait group by 1 as prepare to execute 1 go routine
	crawler.Add()

	// initiate crawling using the fetcher values
	s.ID = NewUlid()
	go crawler.Crawl(s.URL, 25)

	// wait for the initial (and entire) crawl to complete
	crawler.Wait()

	// update crawl details
	s.Duration = time.Now().UTC().Sub(s.ID.Timestamp()).Truncate(time.Second).Seconds()
	s.Relative = crawler.Relative()
	s.Remote = crawler.Remote()

	err := em.Save(s)

	// log the results
	log.Logger.
		Err(err).
		Str("URL", s.URL).
		Int("visited", len(s.Relative)).
		Int("tracked", len(s.Remote)).
		Msg("Sitemap Created")

	return s, err
}

func (s *Sitemap) Delete() (any, error) {
	return nil, s3.New().Delete(s.Key())
}

func (s *Sitemap) FindAll() ([]*Sitemap, error) {
	return em.FindAll(s, regexp.MustCompile(`.*/sitemap/([A-Za-z0-9]{26}/_.json)$`))
}
