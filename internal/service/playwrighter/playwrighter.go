package playwrighter

import (
	"errors"

	"github.com/playwright-community/playwright-go"
)

var (
	defaultBrowserLaunchOptions = playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
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
	}
	defaultPageGotoOptions = playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}
)

type Service interface {
	Browser() playwright.Browser
	Page() (playwright.Page, error)
	Goto(string) (playwright.Page, error)
	Close() error
}

type Client struct {
	player  *playwright.Playwright
	browser playwright.Browser
}

func (c *Client) Browser() playwright.Browser {
	return c.browser
}

func (c *Client) Page() (playwright.Page, error) {
	return c.browser.NewPage()
}

func (c *Client) Goto(url string) (playwright.Page, error) {

	page, err := c.Page()
	if err != nil {
		return nil, err
	}

	var response playwright.Response
	if response, err = page.Goto(url, defaultPageGotoOptions); err != nil {
		return nil, err
	}

	if response.Status() >= 400 {
		return nil, errors.New(response.StatusText())
	}

	return page, nil
}

func (c *Client) Close() error {
	if c.browser == nil {
		return nil
	}
	return c.browser.Close()
}

func New() (Service, error) {

	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}

	var browser playwright.Browser
	if browser, err = pw.Chromium.Launch(); err != nil {
		return nil, err
	}

	return &Client{pw, browser}, nil
}
