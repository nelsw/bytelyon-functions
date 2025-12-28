package model

import (
	"bytelyon-functions/pkg/logger"
	"errors"
	"fmt"
	"log/slog"
	"regexp"

	"github.com/oklog/ulid/v2"
	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

var (
	blockedRegex               = regexp.MustCompile("(google.com/sorry|captcha|unusual traffic)")
	googleSearchInputSelectors = []string{
		"input[name='q']",
		"input[title='Search']",
		"input[aria-label='Search']",
		"textarea[title='Search']",
		"textarea[name='q']",
		"textarea[aria-label='Search']",
		"textarea",
	}
)

func init() {
	if err := playwright.Install(&playwright.RunOptions{
		Logger: slog.New(slogzerolog.Option{
			Level:  slog.LevelDebug,
			Logger: logger.New(),
		}.NewZerologHandler()),
	}); err != nil {
		log.Panic().Err(err).Msg("playwright install")
	}
}

type Prowl struct {
	ID       ulid.ULID
	Headless *bool
	*Prowler
	*playwright.Playwright
	playwright.Browser
	playwright.BrowserContext
}

func (p *Prowl) Close() {
	if err := p.BrowserContext.Close(); err != nil {
		log.Warn().Err(err).Msg("Failed to close Prowl Context")
	} else if err = p.Browser.Close(); err != nil {
		log.Warn().Err(err).Msg("Failed to close Prowl Browser")
	}
}

func (p *Prowl) IsBlocked(ss ...string) error {
	for _, s := range ss {
		if blockedRegex.MatchString(s) {
			return errors.New("blocked: " + s)
		}
	}
	return nil
}

func NewProwl(prowler *Prowler, headless *bool) (p *Prowl, err error) {

	p = &Prowl{
		ID:       NewUlid(),
		Headless: headless,
		Prowler:  prowler,
	}

	if p.Playwright, err = playwright.Run(); err != nil {
		return
	} else if err = p.NewBrowser(); err != nil {
		return
	} else if err = p.NewBrowserContext(); err != nil {
		return
	}
	log.Info().Msg("Prowl - NewProwl")
	return
}

func (p *Prowl) Search() (err error) {

	var page playwright.Page
	var res playwright.Response

	if page, err = p.NewPage(); err != nil {
		return
	} else if res, err = p.GoTo(page, "https://www.google.com"); err != nil {
		return
	} else if err = p.IsBlocked(page.URL(), res.URL()); err != nil {
		return
	} else if err = p.Click(page, googleSearchInputSelectors...); err != nil {
		return
	} else if err = p.Type(page, p.Query); err != nil {
		return
	} else if err = p.Press(page, "Enter"); err != nil {
		return
	} else if err = p.WaitForLoadState(page); err != nil {
		return
	} else if err = p.IsBlocked(page.URL()); err != nil {
		return
	}

	log.Info().Msg("Prowl - Search")

	p.Save(page)

	targetCount := len(p.Targets)
	log.Info().Msgf("Prowl - Targets [%d]", targetCount)

	if targetCount == 0 {
		return
	}

	var locators []playwright.Locator

	if locators, err = page.Locator(fmt.Sprintf(`[data-dtld]`), playwright.PageLocatorOptions{}).All(); err != nil {
		return
	}

	if len(locators) == 0 {
		log.Warn().Msg("Prowl - No Target Locators Found")
	}

	var att string
	for _, l := range locators {

		att, err = l.GetAttribute("data-dtld")
		if err != nil {
			log.Warn().Err(err).Msg("Prowl - Failed to get Target Locator Attribute")
			continue
		}

		log.Debug().Str("found", att).Msg("Prowl - Locator")
		if !p.Targets.Follow(att) {
			continue
		}

		log.Info().Msgf("Prowl - Target Found [%s]", att)

		if page, err = p.NewPage(func() error { return l.Click() }); err == nil {
			//search.SaveState(p.context.StorageState())
			p.Save(page)
			err = errors.Join(page.Close())
		}
	}

	return
}
