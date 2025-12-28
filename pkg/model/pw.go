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

type PW struct {
	*Prowler
	Headless *bool
	*playwright.Playwright
	playwright.Browser
	playwright.BrowserContext
}

func NewPW(prowler *Prowler, headless *bool) (pw *PW, err error) {

	pw = &PW{
		Prowler:  prowler,
		Headless: headless,
	}

	if pw.Playwright, err = playwright.Run(); err != nil {
		return
	} else if err = pw.NewBrowser(); err != nil {
		return
	} else if err = pw.NewBrowserContext(); err != nil {
		return
	}
	log.Info().Msg("PW - NewProwl")
	return
}

func (pw *PW) Close() {
	if err := pw.BrowserContext.Close(); err != nil {
		log.Warn().Err(err).Msg("Failed to close PW Context")
	} else if err = pw.Browser.Close(); err != nil {
		log.Warn().Err(err).Msg("Failed to close PW Browser")
	}
}

func (pw *PW) IsBlocked(ss ...string) error {
	for _, s := range ss {
		if blockedRegex.MatchString(s) {
			return errors.New("blocked: " + s)
		}
	}
	return nil
}

func (pw *PW) Search(prowlID ulid.ULID) (err error) {

	var page playwright.Page
	var res playwright.Response

	if page, err = pw.NewPage(); err != nil {
		return
	} else if res, err = pw.GoTo(page, "https://www.google.com"); err != nil {
		return
	} else if err = pw.IsBlocked(page.URL(), res.URL()); err != nil {
		return
	} else if err = pw.Click(page, googleSearchInputSelectors...); err != nil {
		return
	} else if err = pw.Type(page, pw.Query); err != nil {
		return
	} else if err = pw.Press(page, "Enter"); err != nil {
		return
	} else if err = pw.WaitForLoadState(page); err != nil {
		return
	} else if err = pw.IsBlocked(page.URL()); err != nil {
		return
	}

	log.Info().Msg("PW - Search")

	pw.Save(prowlID, page)

	targetCount := len(pw.Targets)
	log.Info().Msgf("PW - Targets [%d]", targetCount)

	if targetCount == 0 {
		return
	}

	var locators []playwright.Locator

	if locators, err = page.Locator(fmt.Sprintf(`[data-dtld]`), playwright.PageLocatorOptions{}).All(); err != nil {
		return
	}

	if len(locators) == 0 {
		log.Warn().Msg("PW - No Target Locators Found")
	}

	var att string
	for _, l := range locators {

		att, err = l.GetAttribute("data-dtld")
		if err != nil {
			log.Warn().Err(err).Msg("PW - Failed to get Target Locator Attribute")
			continue
		}

		log.Debug().Str("found", att).Msg("PW - Locator")
		if !pw.Targets.Follow(att) {
			continue
		}

		log.Info().Msgf("PW - Target Found [%s]", att)

		if page, err = pw.NewPage(func() error { return l.Click() }); err == nil {
			pw.Save(prowlID, page)
			err = errors.Join(page.Close())
		}
	}

	return
}
