package model

import (
	"bytelyon-functions/pkg/util"
	"regexp"
	"strings"

	"github.com/oklog/ulid/v2"
)

var (
	badExtRegex    = regexp.MustCompile(`^.*\.(jpeg|png|gif|jpg|pdf)$`)
	badAnchorRegex = regexp.MustCompile(`^(#|mailto:|tel:).*`)
)

func (p *Prowler) ProwlSitemap() ulid.ULID {
	prowlID := NewUlid()
	p.Domain = util.Domain(p.URL)
	c := NewCrawler(p)
	c.Add()
	go c.Crawl(p.URL, 15)
	c.Wait()
	p.Relative = c.Relative()
	p.Remote = c.Remote()
	return prowlID
}

func (p *Prowler) Fetch(url string) ([]string, []string, error) {

	if badExtRegex.MatchString(url) {
		return nil, nil, nil
	}

	doc := NewDocument(url)
	if err := doc.Fetch(); err != nil {
		return nil, nil, err
	}

	var relative, remote []string
	for _, a := range doc.anchors() {

		if badAnchorRegex.MatchString(a) {
			continue
		}

		if strings.HasPrefix(a, "?") || strings.HasPrefix(a, "/") {
			relative = append(relative, p.URL+a)
			continue
		}

		u := strings.TrimPrefix(a, "https://")
		u = strings.TrimPrefix(u, "http://")
		u = strings.TrimPrefix(u, "www.")
		u = strings.TrimSuffix(u, "/")
		u = strings.TrimSpace(u)

		if strings.HasPrefix(u, p.Domain) {
			relative = append(relative, a)
		} else {
			remote = append(remote, a)
		}
	}

	return relative, remote, nil
}
