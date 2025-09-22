package model

import (
	"maps"
	"slices"
	"sync"
)

// Crawler encapsulates asynchronous page traversal logic
type Crawler struct {
	Fetcher
	visited map[string]bool
	tracked map[string]bool
	mu      sync.Mutex
	wg      sync.WaitGroup
}

func NewCrawler(fetcher Fetcher) *Crawler {
	return &Crawler{
		Fetcher: fetcher,
		visited: make(map[string]bool),
		tracked: make(map[string]bool),
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
	if depth <= 0 || c.putVisited(URL) {
		return
	}

	// Fetch the url and handle return arguments appropriately
	URLs, links, err := c.Fetch(URL)

	// fail fast on the error, we can't traverse the url
	if err != nil {
		return
	}

	// Store all external links to crawler so that we can make note of egress points
	c.putAllTracked(links)

	// Attempt to crawl each of the domain-specific urls we returned from fetch()
	for _, u := range URLs {
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

func (c *Crawler) putVisited(url string) (ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok = c.visited[url]; !ok {
		c.visited[url] = true
	}
	return ok
}

func (c *Crawler) putAllTracked(urls []string) {
	if len(urls) == 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, u := range urls {
		c.tracked[u] = true
	}
}

func (c *Crawler) Visited() []string {
	return slices.Sorted(maps.Keys(c.visited))
}

func (c *Crawler) Tracked() []string {
	return slices.Sorted(maps.Keys(c.tracked))
}
