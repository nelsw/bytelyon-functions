package model

import (
	"bytelyon-functions/pkg/logger"
	"errors"
	"log/slog"
	"regexp"

	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

var (
	blockedRegex = regexp.MustCompile("(google.com/sorry|captcha|unusual traffic)")
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
	*Prowl
	Headless *bool
	*playwright.Playwright
	playwright.Browser
	playwright.BrowserContext
}

func NewPW(prowl *Prowl, headless *bool) (pw *PW, err error) {

	pw = &PW{
		Prowl:    prowl,
		Headless: headless,
	}

	if pw.Playwright, err = playwright.Run(); err != nil {
		return
	} else if err = pw.NewBrowser(); err != nil {
		return
	} else if err = pw.NewBrowserContext(); err != nil {
		return
	}
	return
}

func (pw *PW) IsBlocked(aa ...any) error {
	for _, a := range aa {
		switch t := a.(type) {
		case playwright.Page:
			if blockedRegex.MatchString(t.URL()) {
				return errors.New("blocked: " + t.URL())
			}
		case playwright.Response:
			if t.Status() >= 400 {
				return errors.New("blocked: " + t.URL())
			}
		}
	}
	return nil
}
