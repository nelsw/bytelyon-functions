package model

import (
	"bytelyon-functions/internal/client/s3"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

type Fetcher interface {
	// Fetch returns the given URL and collects internal urls and external links.
	// Note that we do not crawl external links, but we keep track of them. For reasons.
	Fetch(string) ([]string, []string, error)
}

type Sitemap struct {
	ID      string    `json:"id"`
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Depth   int       `json:"depth"`
	URL     string    `json:"url"`
	Domain  string    `json:"domain"`
	Visited []string  `json:"visited"`
	Tracked []string  `json:"tracked"`
}

func NewSitemap(url string) *Sitemap {
	return &Sitemap{URL: url}
}

func (s *Sitemap) Create(db s3.Client, userID ulid.ULID) (b []byte, err error) {

	// start with the holiest string homogenization ƒ
	s.URL = strings.ToLower(strings.TrimSpace(s.URL))

	// define a new filename friendly ID
	s.ID = base64.URLEncoding.EncodeToString([]byte(s.URL))

	// remove the schemes
	s.Domain = strings.TrimPrefix(s.URL, "http://")
	s.Domain = strings.TrimPrefix(s.Domain, "https://")
	// remove the wild-wild web
	s.Domain = strings.TrimPrefix(s.Domain, "www.")
	// remove the path and query params after the domain extension
	if idx := strings.Index(s.Domain, "/"); idx > 0 {
		s.Domain = s.Domain[:idx]
	}
	// use a loop to remove subdomains ... I think maybe overkill
	for strings.Count(s.Domain, ".") > 1 {
		s.Domain = s.Domain[strings.Index(s.Domain, ".")+1:]
	}

	// use a default depth if invalid value provided
	if s.Depth <= 0 {
		s.Depth = 15
	}

	s.Start = time.Now().UTC()

	// new up a Crawler using a reference to the Sitemap, aka Fetcher
	crawler := NewCrawler(s)

	// increment the crawler wait group by 1 as prepare to execute 1 go routine
	crawler.Add()

	// initiate crawling using the fetcher values
	go crawler.Crawl(s.URL, s.Depth)

	// wait for the initial (and entire) crawl to complete
	crawler.Wait()

	// update crawl details
	s.End = time.Now().UTC()
	s.Visited = crawler.Visited()
	s.Tracked = crawler.Tracked()

	// marshall it
	if b, err = json.Marshal(s); err == nil {
		err = db.Put(UserKey(userID)+"/sitemap/"+s.ID, b)
	}

	// log the results
	log.Logger.
		Err(err).
		Str("ID", s.ID).
		Str("URL", s.URL).
		Int("depth", s.Depth).
		Int("visited", len(s.Visited)).
		Int("tracked", len(s.Tracked)).
		Msg("Sitemap Created")

	return b, err
}

func (s *Sitemap) Fetch(URL string) (urls, links []string, err error) {

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

				if len(href) == 0 {
					continue
				}

				// validate the anchor value and determine where it should be stored
				if string(href[0]) == "/" {
					// check if the anchor is only a path; prefix root URL if so
					urls = append(urls, s.URL+href)
				} else if strings.Contains(href, s.Domain) {
					// append anchor to (inbound) urls and fetch downstream
					urls = append(urls, href)
				} else if !strings.HasPrefix(href, "#") {
					// append anchor to (outbound) links and present as egress point
					links = append(links, href)
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
