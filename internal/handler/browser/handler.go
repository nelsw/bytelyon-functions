package tmp

import (
	"fmt"
	"strings"

	"github.com/playwright-community/playwright-go"
)

// todo - create a layer, in node or w/e to DL browsers at tmp/ms-playwright ... layers _should_ share ephemeral storage

var headless = true
var channel = "chromium"

type Request struct {
	URL string `json:"url"`
}

func Handler(req Request) (err error) {
	fmt.Printf("Received req: %#v\n", req)
	opts := &playwright.RunOptions{
		DriverDirectory:  "/tmp",
		OnlyInstallShell: true,
		Verbose:          true,
		//Browsers: []string{channel},
	}

	var pw *playwright.Playwright
	if pw, err = playwright.Run(opts); err != nil {
		if strings.HasPrefix(err.Error(), "could not get driver instance") {
			// todo install
		}
		if strings.HasPrefix(err.Error(), "please install the driver") {
			err = playwright.Install(opts)
		}

	}
	if err != nil {
		return
	}
	fmt.Println("passed install", pw == nil)
	defer pw.Stop()
	fmt.Println("passed run")
	var br playwright.Browser

	if br, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: &headless,
		//Channel:  &channel,
		Args: []string{
			"--headless=new",
			"--disable-dev-shm-usage",
			"--disable-gpu",
			"--no-sandbox",
			"--remote-debugging-port=9222",
			"--single-process",
			"--window-size=1280x1696",
			"--no-zygote",
			"--enable-javascript",
			"--disable-notifications",
			"--log-level=3",
			"--start-maximized",
			"--enable-automation",
			"--ignore-ssl-errors=yes",
			"--disable-bundled-ppapi-flash",
			"--ignore-certificate-errors",
			"--disable-blink-features=AutomationControlled",
			"--disable-plugins-discovery",
			"--enable-features=NetworkServiceInProcess",
		},
	}); err != nil {
		fmt.Println(err)
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
