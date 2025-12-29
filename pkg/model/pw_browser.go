package model

import (
	. "bytelyon-functions/pkg/util"

	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
)

func (pw *PW) NewBrowser() (err error) {

	pw.Browser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: pw.Headless,
		Timeout:  Ptr(60_000.0),
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
	})

	log.Err(err).Msg("PW - NewBrowser")

	return
}

func (pw *PW) Close() {
	if err := pw.BrowserContext.Close(); err != nil {
		log.Warn().Err(err).Msg("Failed to close PW Context")
	} else if err = pw.Browser.Close(); err != nil {
		log.Warn().Err(err).Msg("Failed to close PW Browser")
	}
}
