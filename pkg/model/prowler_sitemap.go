package model

import (
	"bytelyon-functions/pkg/db"
	"bytelyon-functions/pkg/util"
	"regexp"
	"strings"

	"github.com/oklog/ulid/v2"
	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
)

var (
	badAnchorRegex = regexp.MustCompile(`^(#|mailto:|tel:).*`)
	badExtRegex    = regexp.MustCompile(`^.*\.(jpeg|png|gif|jpg|pdf)$`)
)

type Locator interface {
	Locate(string) ([]string, []string)
}

type ProwlSitemap struct {
	ID       ulid.ULID `json:"id"`
	Domain   string    `json:"domain"`
	Relative []string  `json:"relative"`
	Remote   []string  `json:"remote"`
	*Prowler `json:"-"`
	*PW      `json:"-"`
}

func NewProwlSitemap(p *Prowler) *ProwlSitemap {
	return &ProwlSitemap{
		Prowler: p,
		ID:      NewUlid(),
		Domain:  util.Domain(p.ID),
	}
}

func (p *ProwlSitemap) Dir() string {
	return p.Prowler.Dir()
}

func (p *ProwlSitemap) Key() string {
	return p.Dir() + p.ID.String() + ".json"
}

func (p *ProwlSitemap) Go() ulid.ULID {
	var fn func(bool)
	fn = func(headless bool) {
		var err error
		p.PW, err = NewPW(headless)
		if err != nil {
			log.Warn().Err(err).Msg("ProwlSitemap - Failed to initialize PW")
			return
		}
		defer p.PW.Close()

		log.Info().Msg("ProwlSitemap - Working ... ")

		c := NewProwlSitemapCrawler(p)
		c.Add()
		go c.Crawl(p.Prowler.ID, 3)
		c.Wait()

		p.Relative = c.Relative()
		p.Remote = c.Remote()
	}

	var headless = true
	if fn(true); len(p.Relative) <= 1 {
		log.Debug().Msg("ProwlSitemap - No relative links found with Headless Locator")
		headless = false
		if fn(false); len(p.Relative) <= 1 {
			log.Debug().Msg("ProwlSitemap - No relative links found with Headed Locator")
		}
	}

	log.Info().
		Bool("headless", headless).
		Str("domain", p.Domain).
		Int("relate", len(p.Relative)).
		Int("remote", len(p.Remote)).
		Msg("ProwlSitemap - Done")

	db.Save(p)

	return p.ID
}

func (p *ProwlSitemap) Locate(s string) ([]string, []string) {

	var err error
	var page playwright.Page
	var res playwright.Response

	if page, err = p.PW.NewPage(); err == nil {
		res, err = p.PW.GoTo(page, s)
	}

	if err != nil {
		log.Warn().
			Str("url", s).
			Msg("ProwlSitemap - Failed to locate page")
		return nil, nil
	}

	page.SetDefaultTimeout(0)
	defer page.Close()

	if err = p.PW.IsBlocked(page, res); err != nil {
		log.Warn().Err(err).Msg("ProwlSitemap - Page blocked")
		return nil, nil
	}

	var locators []playwright.Locator
	if locators, err = page.Locator("a").All(); err != nil {
		log.Warn().Err(err).Msg("ProwlSitemap - Failed to locate links")
		return nil, nil
	}

	m := make(map[string]playwright.Locator)
	for _, l := range locators {
		if s, err = l.GetAttribute("href"); err == nil {
			m[s] = l
		}
	}

	var rel, rem []string
	for k := range m {

		k = strings.TrimSuffix(k, "/")

		if badAnchorRegex.MatchString(k) ||
			badExtRegex.MatchString(k) ||
			strings.HasSuffix(k, "@"+p.Domain) {
			continue
		}

		if u := util.Domain(k); u == p.Domain {
			rel = append(rel, k)
			continue
		}

		if strings.HasPrefix(k, "?") || strings.HasPrefix(k, "/") {
			rel = append(rel, p.Prowler.ID+k)
			continue
		}

		rem = append(rem, k)
	}

	return rel, rem
}
