package model

import (
	"bytelyon-functions/pkg/db"
	"bytelyon-functions/pkg/util"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var (
	badAnchorRegex = regexp.MustCompile(`^(#|mailto:|tel:).*`)
	badExtRegex    = regexp.MustCompile(`^.*\.(jpeg|png|gif|jpg|pdf)$`)
)

type ProwlSitemap struct {
	Prowl    *Prowl    `json:"-"`
	ID       ulid.ULID `json:"id"`
	Domain   string    `json:"domain"`
	Relative []string  `json:"relative"`
	Remote   []string  `json:"remote"`
}

func NewProwlSitemap(p *Prowl) *ProwlSitemap {
	return &ProwlSitemap{
		Prowl:  p,
		ID:     NewUlid(),
		Domain: util.Domain(p.Prowler.ID),
	}
}

func (p *ProwlSitemap) String() string {
	return fmt.Sprintf("%s/%s", p.Prowl.Prowler, p.ID)
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

	var fn func(i ...int) (*goquery.Document, error)

	fn = func(i ...int) (*goquery.Document, error) {

		if len(i) == 0 {
			i = append(i, 0)
		}

		if i[0] > 2 {
			return nil, errors.New("Prowl#Sitemap - max fetch attempts reached [3]")
		}
		if i[0] > 0 {
			time.Sleep(time.Second * time.Duration(i[0]) * 2)
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Warn().Err(err).Msgf("Prowl#Sitemap - failed to new up request %s", url)
			return fn(i[0] + 1)
		}
		req.Close = true

		var res *http.Response

		if res, err = http.DefaultClient.Do(req); err != nil {
			log.Warn().
				Err(err).
				Str("url", url).
				Msg("Prowl#Sitemap - Err doing request")
			return fn(i[0] + 1)
		} else if res.StatusCode != http.StatusOK {
			log.Warn().
				Str("url", url).
				Int("status", res.StatusCode).
				Msg("Prowl#Sitemap - Bad Response")
			return fn(i[0] + 1)
		}

		defer res.Body.Close()

		return goquery.NewDocumentFromReader(res.Body)
	}

	doc, err := fn()
	if err != nil {
		return nil, nil, err
	}

	m := make(map[string]bool)
	doc.Find(`a`).Each(func(i int, sel *goquery.Selection) {
		if a, ok := sel.Attr("href"); ok {
			m[strings.TrimSpace(a)] = true
		}
	})

	var relative, remote []string
	for _, a := range slices.Collect(maps.Keys(m)) {

		if badAnchorRegex.MatchString(a) || strings.HasSuffix(a, "@"+p.Domain) {
			continue
		}

		if strings.HasPrefix(a, "?") || strings.HasPrefix(a, "/") {
			relative = append(relative, p.Prowl.Prowler.ID+a)
		} else if u := util.Domain(a); u == p.Domain {
			relative = append(relative, a)
		} else {
			remote = append(remote, a)
		}
	}

	return relative, remote, nil
}
