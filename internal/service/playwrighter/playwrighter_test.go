package playwrighter

import (
	"fmt"
	"log"
	"testing"

	"github.com/playwright-community/playwright-go"
)

func TestNew(t *testing.T) {

	service, err := New()
	if err != nil {
		t.Fatal(err)
	}

	browser := service.Browser()

	var page playwright.Page
	if page, err = browser.NewPage(); err != nil {
		t.Fatal(err)
	}

	url := "https://www.google.com"
	opt := playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}

	var response playwright.Response
	if response, err = page.Goto(url, opt); err != nil {
		t.Fatal(err)
	}

	t.Log(response.Status())
	fmt.Println(page.URL())
	_, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("/tmp/wat.png"),
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGooglyEyes(t *testing.T) {
	service, err := New()
	if err != nil {
		t.Fatal(err)
	}

	browser := service.Browser()

	var page playwright.Page
	if page, err = browser.NewPage(); err != nil {
		t.Fatal(err)
	}

	url := "https://www.google.com"
	opt := playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}

	var response playwright.Response
	if response, err = page.Goto(url, opt); err != nil {
		t.Fatal(err)
	}
	t.Log(response.Status())

	locator := page.GetByTitle("Search")
	if err = locator.Fill("fire blankets"); err != nil {
		t.Fatal(err)
	} else if err = page.GetByRole("button", playwright.PageGetByRoleOptions{
		Name: "Google Search",
	}).Click(); err != nil {
		t.Fatal(err)
	}

	_, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("/tmp/wat.png"),
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestNew2(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch(defaultBrowserLaunchOptions)
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	if _, err = page.Goto("https://google.com"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	locator := page.Locator("textarea[name='q']")
	if err = locator.Fill("fire blankets"); err != nil {
		t.Fatal(err)
	}
	if err = page.GetByRole("button", playwright.PageGetByRoleOptions{
		Name: "Google Search",
	}).Click(); err != nil {
		t.Fatal(err)
	}
	page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{State: playwright.LoadStateLoad})
	_, err = page.Screenshot(playwright.PageScreenshotOptions{
		FullPage: playwright.Bool(true),
		Path:     playwright.String("/Users/connorvanelswyk/GolandProjects/bytelyon-functions/internal/service/playwrighter/tmp/wat.png"),
	})
	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
