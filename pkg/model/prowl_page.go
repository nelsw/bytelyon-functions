package model

import (
	"bytelyon-functions/pkg/service/s3"
	. "bytelyon-functions/pkg/util"
	"encoding/json"
	"errors"
	"fmt"

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

func (p *Prowl) NewPage(ff ...func() error) (page playwright.Page, err error) {

	if len(ff) > 0 {
		page, err = p.BrowserContext.ExpectPage(ff[0], playwright.BrowserContextExpectPageOptions{
			Timeout: playwright.Float(10_000),
		})
	} else if page, err = p.BrowserContext.NewPage(); err == nil {
		err = page.AddInitScript(playwright.Script{Content: Ptr(pageScriptContent)})
	}

	if err != nil {
		log.Warn().Err(err).Msg("Prowl - Failed to NewPage")
	} else {
		log.Info().Str("url", page.URL()).Msg("Prowl - NewPage")
	}

	return
}

func (p *Prowl) GoTo(page playwright.Page, url string) (playwright.Response, error) {

	res, err := page.Goto(url, playwright.PageGotoOptions{
		Timeout:   Ptr(5_000.0),
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})

	if err == nil && !res.Ok() {
		err = errors.New(res.StatusText())
	}

	log.Err(err).Msg("Prowl - GoTo")

	return res, err
}

func (p *Prowl) Click(page playwright.Page, selectors ...string) (err error) {

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
			log.Info().Str("selector", selector).Msg("Prowl - Click")
			return nil
		}

		log.Warn().Err(err).Str("selector", selector).Msg("Prowl - Failed to Click")
	}

	return err
}

func (p *Prowl) WaitForLoadState(page playwright.Page, ls ...playwright.LoadState) error {
	s := playwright.LoadStateNetworkidle
	if len(ls) > 0 {
		s = &ls[0]
	}
	err := page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State:   s,
		Timeout: Ptr(60_000.0),
	})
	log.Err(err).Msg("Prowl - WaitForLoadState")
	return err
}

func (p *Prowl) Save(page playwright.Page) {

	id := NewUlid()
	keyPath := fmt.Sprintf("user/%s/prowler/%s/%s/prowl/%s/page/%s", p.UserID, p.Prowler.Type, p.Prowler.ID, p.ID, id)

	db := s3.New()

	if b, err := page.Screenshot(playwright.PageScreenshotOptions{FullPage: Ptr(true)}); err != nil {
		log.Warn().Err(err).Msg("Prowl - Failed to Screenshot Page")
	} else {
		db.Put(keyPath+"/screenshot.png", b)
	}

	content, err := page.Content()
	if err != nil {
		log.Warn().Err(err).Msg("Prowl - Failed to get Page Content")
	} else {
		db.Put(keyPath+"/content.html", []byte(content))
	}

	var title string
	if title, err = page.Title(); err != nil {
		log.Warn().Err(err).Msg("Prowl - Failed to get Page Title")
	}

	b, _ := json.Marshal(map[string]any{
		"prowler":  p.Prowler,
		"prowl_id": p.ID,
		"id":       id,
		"title":    title,
		"url":      page.URL(),
		"data":     p.Data(page.URL(), content),
	})
	db.Put(keyPath+"/_.json", b)

	log.Info().Stringer("id", id).Msg("Prowl - Saved Page")
}
