package model

import (
	"bytelyon-functions/pkg/db"
	. "bytelyon-functions/pkg/util"
	"encoding/json"
	"errors"

	"github.com/oklog/ulid/v2"
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
			Timeout: playwright.Float(10_000),
		})
	} else if page, err = pw.BrowserContext.NewPage(); err == nil {
		err = page.AddInitScript(playwright.Script{Content: Ptr(pageScriptContent)})
	}

	if err != nil {
		log.Warn().Err(err).Msg("PW - Failed to NewPage")
	} else {
		log.Info().Str("url", page.URL()).Msg("PW - NewPage")
	}

	return
}

func (pw *PW) GoTo(page playwright.Page, url string) (playwright.Response, error) {

	res, err := page.Goto(url, playwright.PageGotoOptions{
		Timeout:   Ptr(5_000.0),
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})

	if err == nil && !res.Ok() {
		err = errors.New(res.StatusText())
	}

	log.Err(err).Msg("PW - GoTo")

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
	log.Err(err).Msg("PW - WaitForLoadState")
	return err
}

func (pw *PW) Save(page playwright.Page) {

	DB := db.NewS3()

	var err error

	var b []byte
	if b, err = page.Screenshot(playwright.PageScreenshotOptions{FullPage: Ptr(true)}); err != nil {
		log.Warn().Err(err).Msg("PW - Failed to Screenshot Page")
	} else {
		DB.Put(pw.Prowler.String()+"/screenshot.png", b)
	}

	var s string
	if s, err = page.Content(); err != nil {
		log.Warn().Err(err).Msg("PW - Failed to get Page Content")
	} else {
		DB.Put(pw.Prowler.String()+"/content.html", []byte(s))
	}

	var p struct {
		ID     ulid.ULID `json:"id"`
		Title  string    `json:"title"`
		URL    string    `json:"url"`
		Domain string    `json:"domain"`
		Data   any       `json:"data"`
	}

	p.ID = NewUlid()
	p.URL = page.URL()
	p.Data = pw.Data(p.URL, s)
	p.Domain = Domain(p.URL)

	if p.Title, err = page.Title(); err != nil {
		log.Warn().Err(err).Msg("PW - Failed to get Page Title")
	}

	if b, err = json.Marshal(p); err != nil {
		log.Warn().Err(err).Msg("PW - Failed to Marshal Page")
	} else {
		DB.Put(pw.Prowler.String()+"/_.json", b)
	}
}
