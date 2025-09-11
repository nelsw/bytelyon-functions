package sitemap

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/s3"
	"bytelyon-functions/internal/model"
	"context"
	"encoding/base64"
	"maps"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

// Crawler encapsulates asynchronous page traversal logic
type Crawler struct {
	Fetcher
	visited map[string]bool
	tracked map[string]bool
	mu      sync.Mutex
	wg      sync.WaitGroup
}

// Crawl is the core function for ... crawling.
// We use sync properties defined in the Crawler to crawl in parallel.
// We also used a couple of maps as a means of bread-crumbing where we've been.
// Ultimately, all we end up doing is logging the results ... for meow 🐱.
func (c *Crawler) Crawl(URL string, depth int) {

	// play it smart and safe - defer done before anything else
	defer c.wg.Done()

	// donny you are beyond your depth. throw some rocks.
	if depth <= 0 {
		return
	}

	// lock the crawler state before accessing the map to avoid a race collission.
	c.mu.Lock()
	// if the url exists in the visited map, we are done with it.
	if _, ok := c.visited[URL]; ok {
		c.mu.Unlock() // don't forget to unlock the state before bailing
		return
	}
	// otherwise we'll be visiting this URL so practively update it's status in the map
	c.visited[URL] = true

	// unlock the crawler state now that we're done with reading and writing the map
	c.mu.Unlock()

	// Fetch the url and handle return arguments appropriately
	URLs, links, err := c.Fetch(URL)

	// Log the fetch results before potentially bailing
	log.Err(err).
		Int("visited", len(URLs)).
		Int("tracked", len(links)).
		Str("URL", URL).
		Msg("Fetch")

	// fail fast on the error; it's logged and Joe is already on it OR HE'S FIRED.
	if err != nil {
		return
	}

	// Store all external links to crawler so that we can make note of egress points
	for _, l := range links {
		c.tracked[l] = true
	}

	// Attempt to crawl each of the domain specific urls we returned from fetch()
	for _, u := range URLs {
		c.wg.Add(1)
		go c.Crawl(u, depth-1)
	}
}

type Fetcher interface {
	Fetch(string) ([]string, []string, error)
}

type Request struct {
	UserID  ulid.ULID `json:"user_id"`
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Depth   int       `json:"depth"`
	URL     string    `json:"url"`
	Domain  string    `json:"domain"`
	Visited []string  `json:"visited"`
	Tracked []string  `json:"tracked"`
}

func (r *Request) Key() string {
	return model.UserKey(r.UserID) + "/sitemap/" + base64.URLEncoding.EncodeToString([]byte(r.URL)) + "/_.json"
}

// Fetch returns the given URL and collects internal urls and external links.
// Note that we do not crawl external links, but we keep track of them. For reasons.
func (r *Request) Fetch(URL string) (urls, links []string, err error) {

	// execute a plain get request and attempt to traverse
	var resp *http.Response
	if resp, err = http.Get(URL); err != nil {
		return
	}
	defer resp.Body.Close()

	// define a document node to do DOM things
	var doc *html.Node
	if doc, err = html.Parse(resp.Body); err != nil {
		return
	}

	// declare the traversal function type here to enable recursion inside the ƒn
	var fn func(*html.Node)

	// define the traversal function, and use comments to love on your future status
	fn = func(n *html.Node) {

		// if the given node is an anchor tag
		if n.Type == html.ElementNode && n.Data == "a" {

			// iterate tag attributes until we find the href
			for _, a := range n.Attr {

				// fail fast if not hyperlink reference
				if a.Key != "href" {
					continue
				}

				// trim it or leave this house right Jeffrey.
				href := strings.TrimSpace(a.Val)

				// validate the anchor value and determine where it should be stored
				if string(href[0]) == "/" {
					// check if the anchor is only a path; prefix root URL if so
					urls = append(urls, r.URL+href)
				} else if strings.Contains(href, r.Domain) {
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

func Handler(ctx context.Context, req Request) {
	Handle(s3.New(ctx), req.UserID, req.URL)
}

func Handle(db s3.Client, userID ulid.ULID, url string) {

	req := Request{
		UserID: userID,
		URL:    url,
		Depth:  15,
	}

	// start with the holiest string homogenization ƒ
	req.URL = strings.ToLower(strings.TrimSpace(req.URL))

	// remove the schemes
	req.Domain = strings.TrimPrefix(req.URL, "http://")
	req.Domain = strings.TrimPrefix(req.Domain, "https://")

	// remove the wild-wild web
	req.Domain = strings.TrimPrefix(req.Domain, "www.")

	// remove the path and query params after the domain extension
	if idx := strings.Index(req.Domain, "/"); idx > 0 {
		req.Domain = req.Domain[:idx]
	}

	// use a loop remove subdomains ... i think maybe overkill
	for strings.Count(req.Domain, ".") > 1 {
		req.Domain = req.Domain[strings.Index(req.Domain, ".")+1:]
	}

	// ... and let's begin
	req.Start = time.Now().UTC()

	// new up a Crawler using a reference to the request, aka Fetcher
	crawler := &Crawler{
		Fetcher: &req,
		visited: make(map[string]bool),
		tracked: make(map[string]bool),
	}

	// increment the crawler wait group by 1 as prepare to execute 1 go routine
	crawler.wg.Add(1)

	// initiate crawling using the fetcher values
	go crawler.Crawl(req.URL, req.Depth)

	// wait for the initial (and entire) crawl to completed
	crawler.wg.Wait()

	// set and end time before saving to s3
	req.End = time.Now().UTC()
	req.Visited = slices.Sorted(maps.Keys(crawler.visited))
	req.Tracked = slices.Sorted(maps.Keys(crawler.tracked))

	// save it
	err := db.Put(req.Key(), app.MustMarshal(req))

	// log the results
	log.Logger.
		Err(err).
		Str("URL", req.URL).
		Int("depth", req.Depth).
		Int("visited", len(req.Visited)).
		Int("tracked", len(req.Tracked)).
		Msg("Fin")
}
