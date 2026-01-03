package model

import (
	"bytelyon-functions/pkg/db"
	"bytelyon-functions/pkg/util"
	"encoding/json"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
)

var (
	googleSearchInputSelectors = []string{
		"input[name='q']",
		"input[title='Search']",
		"input[aria-label='Search']",
		"textarea[title='Search']",
		"textarea[name='q']",
		"textarea[aria-label='Search']",
		"textarea",
	}
)

type ProwlerSearch struct {
	ID       ulid.ULID `json:"id"`
	db.S3    `json:"-"`
	*PW      `json:"-"`
	*Prowler `json:"-"`
}

func NewProwlSearch(p *Prowler) *ProwlerSearch {
	return &ProwlerSearch{
		ID:      NewUlid(),
		S3:      db.NewS3(),
		Prowler: p,
	}
}

func (p *ProwlerSearch) Dir() string {
	return p.Prowler.Dir() + p.ID.String() + "/"
}

func (p *ProwlerSearch) Go() ulid.ULID {
	var fn func(bool) ulid.ULID

	fn = func(headless bool) ulid.ULID {

		var err error
		p.PW, err = NewPW(headless)
		if err != nil {
			log.Warn().Err(err).Msg("ProwlerSearch - Failed to initialize PW")
			return ulid.Zero
		}
		defer p.PW.Close()

		log.Info().Bool("headless", headless).Msg("ProwlerSearch - Working ... ")

		var prowled ulid.ULID
		if prowled, err = p.worker(); err != nil && headless {
			log.Warn().Err(err).Msg("ProwlerSearch - Headless Failed; retrying with head ...")
			return fn(false)
		}

		if err != nil {
			log.Warn().Err(err).Bool("headless", headless).Msg("ProwlerSearch - Failed!")
			return ulid.Zero
		}

		log.Info().Bool("headless", headless).Msg("ProwlerSearch - Success!")
		return prowled
	}

	return fn(true)
}

func (p *ProwlerSearch) worker() (prowled ulid.ULID, err error) {

	defer p.PW.Close()

	var page playwright.Page
	if page, err = p.PW.NewPage(); err != nil {
		return
	}
	defer page.Close()

	var res playwright.Response
	if res, err = p.PW.GoTo(page, "https://www.google.com"); err != nil {
		return
	} else if err = p.PW.IsBlocked(page, res); err != nil {
		return
	} else if err = p.PW.Click(page, googleSearchInputSelectors...); err != nil {
		return
	} else if err = p.PW.Type(page, p.Prowler.ID); err != nil {
		return
	} else if err = p.PW.Press(page, "Enter"); err != nil {
		return
	} else if err = p.PW.WaitForLoadState(page); err != nil {
		return
	} else if err = p.PW.IsBlocked(page); err != nil {
		return
	}

	prowled = p.save(page)

	targetCount := len(p.Targets)
	log.Info().Msgf("ProwlerSearch - Pages [%d]", targetCount)

	if targetCount == 0 {
		return
	}

	var locators []playwright.Locator
	if locators, err = page.Locator(fmt.Sprintf(`[data-dtld]`), playwright.PageLocatorOptions{}).All(); err != nil {
		return
	}

	if len(locators) == 0 {
		log.Warn().Msg("ProwlerSearch - No Target Locators Found")
		return
	}

	var att string
	for _, l := range locators {

		if att, err = l.GetAttribute("data-dtld", playwright.LocatorGetAttributeOptions{
			Timeout: util.Ptr(5_000.0),
		}); err != nil {
			log.Warn().Err(err).Msg("ProwlerSearch - Failed to get Target Locator Attribute")
			continue
		}

		log.Debug().Str("found", att).Msg("ProwlerSearch - Locator")
		if !p.Targets.Follow(att) {
			continue
		}

		log.Info().Msgf("ProwlerSearch - Target Found [%s]", att)

		if err = l.Click(playwright.LocatorClickOptions{Timeout: util.Ptr(5_000.0)}); err != nil {
			log.Warn().Err(err).Msg("ProwlerSearch - Failed to Click Target Locator")
			continue
		}

		targetPage, pageErr := p.PW.NewPage(func() error { return l.Click() })
		if pageErr != nil {
			log.Warn().Err(pageErr).Msg("ProwlerSearch - Failed to Click Target")
		} else {
			prowled = p.save(targetPage)
			targetPage.Close()
		}
	}

	return prowled, nil
}

func (p *ProwlerSearch) save(page playwright.Page) (prowled ulid.ULID) {

	prowled = NewUlid()
	url := page.URL()
	domain := util.Domain(url)

	path := p.Dir()
	if domain == "google.com" {
		path += "serp/"
	} else {
		path += "target/"
	}
	path += prowled.String()

	var err error

	var img []byte
	if img, err = page.Screenshot(playwright.PageScreenshotOptions{FullPage: util.Ptr(true)}); err != nil {
		log.Warn().Err(err).Str("domain", domain).Msg("PW - Failed to Screenshot Page")
	} else {
		p.S3.Put(path+".png", img)
	}

	var content string
	if content, err = page.Content(); err != nil {
		log.Warn().Err(err).Str("domain", domain).Msg("PW - Failed to get Page Content")
	} else {
		p.S3.Put(path+".html", []byte(content))
	}

	var title string
	if title, err = page.Title(); err != nil {
		log.Warn().Err(err).Str("domain", domain).Msg("PW - Failed to get Page Title")
	}

	m := map[string]any{
		"url":    url,
		"domain": domain,
		"title":  title,
	}

	if m["domain"] == "google.com" {
		m["results"] = p.PW.Data(url, content)
		m["targets"] = p.Targets
	}

	var data []byte
	if data, err = json.Marshal(&m); err != nil {
		log.Warn().Err(err).Str("domain", domain).Msg("PW - Failed to Marshal Page")
	} else {
		p.S3.Put(path+".json", data)
	}

	log.Info().
		Str("domain", domain).
		Msg("PW - Saved Page")

	return
}
