package prowl

import (
	. "bytelyon-functions/pkg/util"

	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
)

func (p *Prowler) NewBrowser() (err error) {
	opts := playwright.BrowserTypeLaunchOptions{
		Headless: &p.headless,
		Timeout:  Ptr(2 * 60_000.0),
		Args: []string{
			"--disable-accelerated-2d-canvas",
			"--disable-background-networking",
			"--disable-background-timer-throttling",
			"--disable-backgrounding-occluded-windows",
			"--disable-blink-features=AutomationControlled",
			"--disable-breakpad",
			"--disable-component-extensions-with-background-pages",
			"--disable-dev-shm-usage",
			"--disable-extensions",
			"--disable-features=IsolateOrigins,site-per-process",
			"--disable-features=TranslateUI",
			"--disable-gpu",
			"--disable-ipc-flooding-protection",
			"--disable-renderer-backgrounding",
			"--disable-setuid-sandbox",
			"--disable-site-isolation-trials",
			"--disable-web-security",
			"--enable-features=NetworkService,NetworkServiceInProcess",
			"--force-color-profile=srgb",
			"--hide-scrollbars",
			"--metrics-recording-only",
			"--mute-audio",
			"--no-first-run",
			"--no-sandbox",
			"--no-zygote",
		},
		IgnoreDefaultArgs: []string{
			"--enable-automation",
		},
	}
	switch p.BrowserType {
	case Chromium:
		p.Browser, err = p.Chromium.Launch(opts)
	case Firefox:
		p.Browser, err = p.Firefox.Launch(opts)
	case WebKit:
		p.Browser, err = p.WebKit.Launch(opts)
	}

	log.Err(err).Msg("Prowler - NewBrowser")

	return
}
