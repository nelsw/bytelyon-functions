package tmp

import (
	"fmt"
	"strings"

	"github.com/playwright-community/playwright-go"
)

// todo - create a layer, in node or w/e to DL browsers at tmp/ms-playwright ... layers _should_ share ephemeral storage

var headless = false
var channel = "chrome"

type Request struct {
	URL string `json:"url"`
}

func Handler(req Request) (err error) {

	opts := &playwright.RunOptions{
		DriverDirectory:     "/tmp/ms-playwright",
		OnlyInstallShell:    false,
		SkipInstallBrowsers: false,
		DryRun:              false,
	}

	var pw *playwright.Playwright
	if pw, err = playwright.Run(opts); err != nil {
		if strings.HasPrefix(err.Error(), "could not get driver instance") {
			// todo install
		}
		if strings.HasPrefix(err.Error(), "please install the driver") {
			err = playwright.Install(opts)
		}
		if err != nil {
			return
		}
	}
	defer pw.Stop()

	var br playwright.Browser
	if br, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Channel:  &channel,
		Headless: &headless,
		Args: []string{
			"--disable-blink-features",
			"--single-process",
			"--no-zygote",
			"--no-sandbox",
		},
	}); err != nil {
		return
	}
	defer br.Close()

	var bc playwright.BrowserContext
	if bc, err = br.NewContext(playwright.BrowserNewContextOptions{}); err != nil {
		return
	}
	defer bc.Close()

	var pg playwright.Page
	if pg, err = bc.NewPage(); err != nil {
		return
	}
	defer pg.Close()

	if _, err = pg.Goto(req.URL); err != nil {
		return
	}

	fmt.Println(pg.Content())

	return
}
