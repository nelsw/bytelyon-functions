package pw

import (
	"bytelyon-functions/pkg/util/ptr"

	"github.com/playwright-community/playwright-go"
)

type Service interface {
	Search(string) (string, []byte, error)
}

func New() (Service, error) {
	p, err := playwright.Run()
	if err != nil {
		return nil, err
	}

	var proxy *playwright.Proxy
	if proxy, err = NewProxy(); err != nil {
		return nil, err
	}

	var b playwright.Browser
	b, err = p.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: ptr.True(),
		Timeout:  ptr.Float64(2 * 60_000),
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
		Proxy: proxy,
	})
	if err != nil {
		return nil, err
	}

	return &Browser{b, proxy}, nil
}
