package model

import (
	"bytelyon-functions/pkg/db"
	"bytelyon-functions/pkg/util"
	"regexp"
	"strings"
)

var (
	badAnchorRegex = regexp.MustCompile(`^(#|mailto:|tel:).*`)
	badExtRegex    = regexp.MustCompile(`^.*\.(jpeg|png|gif|jpg|pdf)$`)
)

type ProwlSitemap struct {
	Prowl    *Prowl   `json:"prowl"`
	Domain   string   `json:"domain"`
	Relative []string `json:"relative"`
	Remote   []string `json:"remote"`
}

func NewProwlSitemap(p *Prowl) *ProwlSitemap {
	return &ProwlSitemap{
		Prowl:  p,
		Domain: util.Domain(p.Prowler.ID),
	}
}

func (p *ProwlSitemap) String() string {
	return p.Prowl.String()
}

func (p *ProwlSitemap) Go() {

	c := NewCrawler(p)
	c.Add()
	go c.Crawl(p.Prowl.Prowler.ID, 15)
	c.Wait()

	p.Relative = c.Relative()
	p.Remote = c.Remote()

	db.Save(p)
}

func (p *ProwlSitemap) Fetch(url string) ([]string, []string, error) {

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
			relative = append(relative, p.Prowl.Prowler.ID+a)
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
