package model

import (
	. "bytelyon-functions/pkg/util"
	"errors"

	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
)

const (
	pageScriptContent = `() => {
  Object.defineProperty(window.screen, "width", { get: () => 1920 });
  Object.defineProperty(window.screen, "height", { get: () => 1080 });
  Object.defineProperty(window.screen, "colorDepth", { get: () => 24 });
  Object.defineProperty(window.screen, "pixelDepth", { get: () => 24 });
}`
)

func (pw *PW) NewPage(ff ...func() error) (page playwright.Page, err error) {

	if len(ff) > 0 {
		page, err = pw.BrowserContext.ExpectPage(ff[0], playwright.BrowserContextExpectPageOptions{
			Timeout: playwright.Float(30_000),
		})
	} else if page, err = pw.BrowserContext.NewPage(); err == nil {
		page.SetDefaultTimeout(30_000)
		err = page.AddInitScript(playwright.Script{Content: Ptr(pageScriptContent)})
	}

	if err != nil {
		log.Warn().Err(err).Msg("PW - Failed to NewPage")
	}

	return
}

func (pw *PW) GoTo(page playwright.Page, url string) (playwright.Response, error) {

	res, err := page.Goto(url, playwright.PageGotoOptions{
		Timeout:   Ptr(30_000.0),
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})

	if err == nil && !res.Ok() {
		err = errors.New(res.StatusText())
	}

	log.Err(err).Str("url", url).Msg("PW - GoTo")

	return res, err
}

func (pw *PW) Click(page playwright.Page, selectors ...string) (err error) {

	var locator playwright.Locator
	for _, selector := range selectors {

		if locator = page.Locator(selector); locator == nil {
			continue
		}

		var n int
		if n, err = locator.Count(); n == 0 {
			continue
		}

		if err = locator.Click(playwright.LocatorClickOptions{Delay: Ptr(Between(200, 500.0))}); err == nil {
			log.Info().Str("selector", selector).Msg("PW - Click")
			return nil
		}

		log.Warn().Err(err).Str("selector", selector).Msg("PW - Failed to Click")
	}

	return err
}

func (pw *PW) WaitForLoadState(page playwright.Page, ls ...playwright.LoadState) error {
	s := playwright.LoadStateNetworkidle
	if len(ls) > 0 {
		s = &ls[0]
	}
	err := page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State:   s,
		Timeout: Ptr(60_000.0),
	})
	if err != nil {
		log.Err(err).Msg("PW - WaitForLoadState")
	}
	return err
}
