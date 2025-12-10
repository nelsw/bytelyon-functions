package model

import (
	"maps"
	"slices"
	"sync"

	"github.com/rs/zerolog/log"
)

// Crawler encapsulates asynchronous page traversal logic
type Crawler struct {
	Fetcher
	relative map[string]string
	remote   map[string]string
	mu       sync.Mutex
	wg       sync.WaitGroup
}

func NewCrawler(fetcher Fetcher) *Crawler {
	return &Crawler{
		Fetcher:  fetcher,
		relative: make(map[string]string),
		remote:   make(map[string]string),
	}
}

// Crawl is the core function for ... crawling.
// We use sync properties defined in the Crawler to crawl in parallel.
// We also used a couple of maps as a means of bread-crumbing where we've been.
// Ultimately, all we end up doing is logging the results ... for meow 🐱.
func (c *Crawler) Crawl(URL string, depth int) {

	// play it smart and safe - defer done before anything else
	defer c.wg.Done()

	// fail fast if we're past our depth or if we've already visited the URL
	if depth <= 0 || c.putRelative(URL) {
		return
	}

	// Fetch the url and handle return arguments appropriately
	relativeURLs, remoteURLs, err := c.Fetch(URL)

	// Store all external links to crawler so that we can make note of egress points
	c.putAllRemote(remoteURLs)

	// fail fast on the error, no urls or links to follow
	if err != nil {
		log.Err(err).Str("URL", URL).Msg("failed to fetch")
		return
	}

	// Attempt to crawl each of the domain-specific urls we returned from fetch()
	for _, u := range relativeURLs {
		c.Add()
		go c.Crawl(u, depth-1)
	}
}

func (c *Crawler) Add() {
	c.wg.Add(1)
}

func (c *Crawler) Wait() {
	c.wg.Wait()
}

func (c *Crawler) putRelative(url string) (ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok = c.relative[url]; !ok {
		c.relative[url] = url
	}
	return ok
}

func (c *Crawler) putAllRemote(urls []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, url := range urls {
		c.remote[url] = url
	}
}

func (c *Crawler) Relative() []string {
	return slices.Sorted(maps.Keys(c.relative))
}

func (c *Crawler) Remote() []string {
	return slices.Sorted(maps.Keys(c.remote))
}
