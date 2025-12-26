package prowl

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

func (p *Prowler) NewPage() (err error) {
	if p.Page, err = p.BrowserContext.NewPage(); err == nil {
		err = p.Page.AddInitScript(playwright.Script{Content: Ptr(pageScriptContent)})
	}
	log.Err(err).Msg("Prowler - NewPage")
	return
}

func (p *Prowler) GoTo(url string) (playwright.Response, error) {

	res, err := p.Page.Goto(url, playwright.PageGotoOptions{
		Timeout:   Ptr(5_000.0),
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})

	if err == nil && !res.Ok() {
		err = errors.New(res.StatusText())
	}

	log.Err(err).Msg("Prowler - GoTo")

	return res, err
}

func (p *Prowler) Click(selectors ...string) error {

	var err error
	for _, selector := range selectors {

		locator := p.Page.Locator(selector)
		if locator == nil {
			continue
		}

		var n int
		if n, err = locator.Count(); n == 0 {
			continue
		}
		err = locator.Click(playwright.LocatorClickOptions{
			Delay: Ptr(Between(200, 500.0)),
		})
		break
	}
	log.Err(err).Bool("headless", p.headless).Strs("selectors", selectors).Msg("Click")
	return err
}

func (p *Prowler) WaitForLoadState(ls ...playwright.LoadState) error {
	s := playwright.LoadStateNetworkidle
	if len(ls) > 0 {
		s = &ls[0]
	}
	err := p.Page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State:   s,
		Timeout: Ptr(60_000.0),
	})
	log.Err(err).Msg("Prowler - WaitForLoadState")
	return err
}

func (p *Prowler) Title() string {

	title, err := p.BrowserContext.Pages()[len(p.BrowserContext.Pages())-1].Title()

	if err != nil {
		log.Warn().Err(err).Bool("headless", p.headless).Msg("Title")
		return ""
	}

	log.Info().Str("title", title).Bool("headless", p.headless).Msg("Title")
	return title
}

func (p *Prowler) Screenshot() []byte {
	b, err := p.BrowserContext.Pages()[len(p.BrowserContext.Pages())-1].Screenshot(playwright.PageScreenshotOptions{FullPage: playwright.Bool(true)})
	log.Err(err).
		Bool("headless", p.headless).
		Int("size", len(b)).
		Msg("Screenshot")
	return b
}

func (p *Prowler) Content() string {
	content, err := p.BrowserContext.Pages()[len(p.BrowserContext.Pages())-1].Content()
	log.Err(err).
		Bool("headless", p.headless).
		Int("size", len(content)).
		Msg("Content")
	return content
}
