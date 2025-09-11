package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type Fetcher interface {
	Domain() string
	URL() string
	Fetch(string) ([]string, error) // Returns extracted URLs
}

// SimpleFetcher implements Fetcher
type SimpleFetcher struct {
	domain string
	rawUrl string
	depth  int
}

func MakeSimpleFetcher(u string, depth int) SimpleFetcher {

	// remove the schemes
	d := strings.TrimPrefix(u, "http://")
	d = strings.TrimPrefix(d, "https://")

	// remove the wild-wild web
	d = strings.TrimPrefix(d, "www.")

	// remove the path and query params after the domain extension
	if idx := strings.Index(d, "/"); idx > 0 {
		d = d[:idx]
	}

	// remove subdomains
	for strings.Count(d, ".") > 1 {
		d = d[strings.Index(d, ".")+1:]
	}

	// like all things str matching ... lowercase it
	d = strings.ToLower(d)

	return SimpleFetcher{
		domain: d,
		rawUrl: u,
		depth:  depth,
	}
}

func (s SimpleFetcher) Domain() string {
	return s.domain
}

func (s SimpleFetcher) URL() string {
	return s.rawUrl
}

func (s SimpleFetcher) Fetch(url string) (urls []string, err error) {

	var resp *http.Response
	if resp, err = http.Get(url); err != nil {
		return
	}
	defer resp.Body.Close()

	var doc *html.Node
	if doc, err = html.Parse(resp.Body); err != nil {
		return
	}

	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					urls = append(urls, a.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return
}

// Crawler manages the crawling process
type Crawler struct {
	visited map[string]bool
	mu      sync.Mutex
	wg      sync.WaitGroup
	fetcher Fetcher
}

func NewCrawler(fetcher Fetcher) *Crawler {
	return &Crawler{
		visited: make(map[string]bool),
		fetcher: fetcher,
	}
}

func (c *Crawler) Crawl(url string, depth int) {
	defer c.wg.Done()

	if depth <= 0 {
		return
	}

	c.mu.Lock()
	if _, ok := c.visited[url]; ok {
		c.mu.Unlock()
		return
	}
	c.visited[url] = true
	c.mu.Unlock()

	fmt.Printf("Crawling: %s (Depth: %d)\n", url, depth)

	links, err := c.fetcher.Fetch(url)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", url, err)
		return
	}

	for _, link := range links {
		if string(link[0]) == "/" {
			link = c.fetcher.URL() + link
		} else if !strings.Contains(link, c.fetcher.Domain()) {
			continue
		}
		c.wg.Add(1)
		go c.Crawl(link, depth-1)
	}
}

func main() {
	fetcher := MakeSimpleFetcher("https://account.li-fire.com", 15)

	crawler := NewCrawler(fetcher)

	crawler.wg.Add(1)
	go crawler.Crawl(fetcher.rawUrl, fetcher.depth)
	crawler.wg.Wait()

	fmt.Println("Crawling finished. ", len(crawler.visited))
}
