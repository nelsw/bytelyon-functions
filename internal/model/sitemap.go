package model

import (
	"bytelyon-functions/internal/client/s3"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

type Fetcher interface {
	// Fetch returns the given URL and collects internal urls and external links.
	// Note that we do not crawl external links, but we keep track of them. For reasons.
	Fetch(string) ([]string, []string, error)
}
type Sitemaps []Sitemap
type Sitemap struct {
	User     *User     `json:"-"`
	ID       string    `json:"id"`
	URL      string    `json:"url"`
	Domain   string    `json:"domain"`
	Start    time.Time `json:"start"`
	Duration float64   `json:"duration"`
	Relative []string  `json:"relative"`
	Remote   []string  `json:"remote"`
}

func (s *Sitemap) Path() string {
	return s.User.Path() + "/sitemap"
}

func (s *Sitemap) Key() string {
	return s.Path() + "/" + s.ID + "/_.json"
}

func NewSitemap(req events.APIGatewayV2HTTPRequest) (*Sitemap, error) {

	u, err := NewUser(req)
	if err != nil {
		return nil, err
	}

	encodedURL := req.QueryStringParameters["url"]
	if encodedURL == "" {
		return &Sitemap{
			User: u,
		}, nil
	}

	var b []byte
	if b, err = base64.URLEncoding.DecodeString(encodedURL); err != nil {
		return nil, err
	}

	decodedURL := string(b)
	if _, err = url.ParseRequestURI(decodedURL); err != nil {
		return nil, err
	}

	return &Sitemap{
		ID:   encodedURL,
		User: u,
		URL:  decodedURL,
	}, nil
}

func (s *Sitemap) Fetch(URL string) (relative, remote []string, err error) {

	// filter out mailto: and email addresses
	if strings.HasPrefix(URL, "mailto:") || strings.HasSuffix(URL, "@"+s.Domain) {
		return
	}

	// execute a plain get request and attempt to traverse
	var resp *http.Response
	if resp, err = s.get(URL, 0); err != nil {
		return
	}

	// define a document node to do DOM things
	var doc *html.Node
	doc, err = html.Parse(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
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

				// trim it or leave this house right, Jeffrey.
				href := strings.TrimSpace(a.Val)

				// trim potential trailing slashes
				href = strings.TrimSuffix(href, "/")

				// validate the anchor value
				if len(href) == 0 {
					continue
				}

				// check if the anchor is valid and not a hash/mail/tel reference or email address
				if strings.HasSuffix(href, "@"+s.Domain) ||
					strings.HasPrefix(href, "mailto:") ||
					strings.HasPrefix(href, "tel:") ||
					strings.HasPrefix(href, "#") {
					continue
				}

				// check if the anchor is only a path; prefix root URL if so
				if v := string(href[0]); v == "/" || v == "?" {
					relative = append(relative, s.URL+href)
					continue
				}

				// check if the anchor starts with the domain + "https://www."
				if i := strings.Index(href, s.Domain); i >= 0 && i <= 12 {
					relative = append(relative, href)
					continue
				}

				// check if the anchor is a valid remote URL
				if _, err = url.Parse(href); err == nil {
					remote = append(remote, href)
				}
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

func (s *Sitemap) get(URL string, attempt int) (*http.Response, error) {

	if attempt > 0 {
		time.Sleep(time.Duration(attempt*3) * time.Second)
	}

	if resp, err := http.Get(URL); err == nil {
		return resp, nil
	} else if attempt <= 3 {
		return s.get(URL, attempt+1)
	} else {
		log.Err(err).Str("URL", URL).Msg("failed to http.Get")
		return nil, err
	}
}

func (s *Sitemap) Create() (*Sitemap, error) {

	// remove the schemes
	s.Domain = strings.TrimPrefix(s.URL, "http://")
	s.Domain = strings.TrimPrefix(s.Domain, "https://")
	// remove the wild-wild web
	s.Domain = strings.TrimPrefix(s.Domain, "www.")
	// remove the path and query params after the domain extension
	if i := strings.Index(s.Domain, "/"); i > 0 {
		s.Domain = s.Domain[:i]
	}

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
		Str("URL", s.URL).
		Int("visited", len(s.Relative)).
		Int("tracked", len(s.Remote)).
		Msg("Sitemap Created")

	return s, err
}

func (s *Sitemap) FindAll() (Sitemaps, error) {

	db := s3.New()

	keys, err := db.Keys(s.Path(), "", 1000)
	if err != nil {
		return nil, err
	}

	var vv Sitemaps
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
