package prowl

import (
	"bytelyon-functions/pkg/logger"
	"bytelyon-functions/pkg/model"
	"errors"
	"log/slog"
	"regexp"

	"github.com/oklog/ulid/v2"
	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

var (
	runOpts = playwright.RunOptions{
		Logger: slog.New(slogzerolog.Option{
			Level:  slog.LevelDebug,
			Logger: logger.New(),
		}.NewZerologHandler()),
	}
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

type Prowler struct {
	userID   ulid.ULID
	headless bool
	BrowserType
	*playwright.Playwright
	playwright.Browser
	playwright.BrowserContext
	playwright.Page
}

func (p *Prowler) Close() {
	if err := p.BrowserContext.Close(); err != nil {
		log.Warn().Err(err).Msg("Failed to close Prowler Context")
	} else if err = p.Browser.Close(); err != nil {
		log.Warn().Err(err).Msg("Failed to close Prowler Browser")
	}
}

func (p *Prowler) IsBlocked(ss ...string) error {
	for _, s := range ss {
		if blockedRegex.MatchString(s) {
			return errors.New("blocked: " + s)
		}
	}
	return nil
}

func New(a ...any) (p *Prowler, err error) {

	p = &Prowler{BrowserType: RandomBrowserType()}
	for _, v := range a {
		switch v.(type) {
		case model.User:
			p.userID = v.(model.User).ID
		case ulid.ULID:
			p.userID = v.(ulid.ULID)
		case bool:
			p.headless = v.(bool)
		}
	}

	err = playwright.Install(&runOpts)
	log.Err(err).Msg("Prowler Install")
	if err != nil {
		return nil, err
	}

	p.Playwright, err = playwright.Run()
	log.Err(err).Msg("Prowler - Run")

	if err = p.NewBrowser(); err != nil {
		return
	} else if err = p.NewBrowserContext(); err != nil {
		return
	} else if err = p.NewPage(); err != nil {
		return
	}
	log.Info().Msg("Prowler Ready")
	return
}

func (p *Prowler) Hunt() error {

	return nil
}

func (p *Prowler) search(query string) (err error) {

	log.Info().Msg("Prowler#Search - ...")

	var res playwright.Response
	if res, err = p.GoTo("https://www.google.com"); err != nil {
		return
	} else if err = p.IsBlocked(p.Page.URL(), res.URL()); err != nil {
		return
	} else if err = p.Click(googleSearchInputSelectors...); err != nil {
		return
	} else if err = p.Type(query); err != nil {
		return
	} else if err = p.Press("Enter"); err != nil {
		return
	} else if err = p.IsBlocked(p.Page.URL()); err != nil {
		return
	}

	log.Info().Msg("Prowler#Search - Successful")

	return
}
