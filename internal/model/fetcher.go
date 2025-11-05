package model

type Fetcher interface {
	// Fetch returns the given URL and collects internal urls and external links.
	// Note that we do not crawl external links, but we keep track of them. For reasons.
	Fetch(URL) ([]URL, []URL, error)
}
