package model

import (
	"maps"
	"slices"
	"sync"
)

type ProwlerSitemapper struct {
	Locator
	relative map[string]bool
	remote   map[string]bool
	mu       sync.Mutex
	wg       sync.WaitGroup
}

func NewProwlSitemapCrawler(p *ProwlSitemap) *ProwlerSitemapper {
	return &ProwlerSitemapper{
		Locator:  p,
		relative: make(map[string]bool),
		remote:   make(map[string]bool),
	}
}

func (c *ProwlerSitemapper) Crawl(url string, depth int) {

	defer c.wg.Done()

	if depth <= 0 || c.putRelative(url) {
		return
	}

	rel, rem := c.Locate(url)

	c.putAllRemote(rem)

	for _, u := range rel {
		c.Add()
		go c.Crawl(u, depth-1)
	}
}

func (c *ProwlerSitemapper) Add() {
	c.wg.Add(1)
}

func (c *ProwlerSitemapper) Wait() {
	c.wg.Wait()
}

func (c *ProwlerSitemapper) putRelative(s string) (ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok = c.relative[s]; !ok {
		c.relative[s] = true
	}
	return ok
}

func (c *ProwlerSitemapper) putAllRemote(urls []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, url := range urls {
		c.remote[url] = true
	}
}

func (c *ProwlerSitemapper) Relative() []string {
	return slices.Sorted(maps.Keys(c.relative))
}

func (c *ProwlerSitemapper) Remote() []string {
	return slices.Sorted(maps.Keys(c.remote))
}
