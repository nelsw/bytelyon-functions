package model

import (
	"bytelyon-functions/pkg/service/s3"
	"encoding/json"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

type Sitemap struct {
	User     *User     `json:"-"`
	ID       string    `json:"id"`
	URL      URL       `json:"url"`
	Domain   string    `json:"domain"`
	Start    time.Time `json:"start"`
	Duration float64   `json:"duration"`
	Relative []string  `json:"relative"`
	Remote   []string  `json:"remote"`
}

func (s *Sitemap) Path() string {
	return s.User.Dir() + "/sitemap"
}

func (s *Sitemap) Key() string {
	return s.Path() + "/" + s.ID + "/_.json"
}

func NewSitemap(user *User, encodedURL string) (*Sitemap, error) {

	url, err := DecodeURL(encodedURL)
	if err != nil {
		return nil, err
	}

	return &Sitemap{
		ID:     encodedURL,
		User:   user,
		URL:    url,
		Domain: url.Domain(),
	}, nil
}

func (s *Sitemap) Fetch(u URL) (relative, remote []URL, err error) {

	var doc *html.Node
	if doc, err = u.Document(); err != nil {
		return
	}

	// declare the traversal function type here to enable recursion inside the ƒn
	var fn func(*html.Node)

	// define the traversal function and use comments to love on your future status
	fn = func(n *html.Node) {

		// if the given node is an anchor tag
		if n.Type == html.ElementNode && n.Data == "a" {

			// iterate tag attributes until we find the href
			for _, a := range n.Attr {

				// fail fast if not hyperlink reference
				if a.Key != "href" {
					continue
				}

				href := MakeURL(a.Val)

				if err = href.Validate(); err != nil {
					continue
				}

				// check if the anchor is valid and not a hash/mail/tel reference or email address
				if href.StartsWith("mailto:", "tel:", "#") || href.EndsWith("@"+s.Domain) {
					continue
				}

				// check if the anchor is only a path; prefix root URL if so
				if href.StartsWith("/", "?") {
					relative = append(relative, s.URL.Append(href))
					continue
				}

				// check if the anchor starts with the domain + "https://www."
				if href.Domain() == s.Domain {
					relative = append(relative, href)
					continue
				}

				remote = append(remote, href)
			}
		}

		// continue traversing every sibling per child. give em noogies.
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fn(c)
		}
	}

	// and ... all together ... RECURSE!
	fn(doc)

	return
}

func (s *Sitemap) Create() (*Sitemap, error) {

	// new up a Crawler using a reference to the Sitemap, aka Fetcher
	crawler := NewCrawler(s)

	// increment the crawler wait group by 1 as prepare to execute 1 go routine
	crawler.Add()

	// initiate crawling using the fetcher values
	s.Start = time.Now().UTC()
	go crawler.Crawl(s.URL, 25)

	// wait for the initial (and entire) crawl to complete
	crawler.Wait()

	// update crawl details
	s.Duration = time.Now().UTC().Sub(s.Start).Truncate(time.Second).Seconds()
	s.Relative = crawler.Relative()
	s.Remote = crawler.Remote()

	// marshall it
	b, err := json.Marshal(s)
	if err == nil {
		err = s3.New().Put(s.Key(), b)
	}

	// log the results
	log.Logger.
		Err(err).
		Str("URL", s.URL.String()).
		Int("visited", len(s.Relative)).
		Int("tracked", len(s.Remote)).
		Msg("Sitemap Created")

	return s, err
}

func (s *Sitemap) FindAll() ([]Sitemap, error) {

	db := s3.New()

	keys, err := db.Keys(s.Path(), "", 1000)
	if err != nil {
		return nil, err
	}

	var vv []Sitemap
	for _, k := range keys {

		o, e := db.Get(k)
		if e != nil {
			err = errors.Join(err, e)
			continue
		}

		var v Sitemap
		if e = json.Unmarshal(o, &v); e != nil {
			err = errors.Join(err, e)
			continue
		}

		vv = append(vv, v)
	}

	return vv, err
}

func (s *Sitemap) Delete() (any, error) {
	return nil, s3.New().Delete(s.Key())
}
