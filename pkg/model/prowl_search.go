package model

import (
	"errors"
	"fmt"

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

type ProwlSearch struct {
	*Prowl
}

func NewProwlSearch(p *Prowl) *ProwlSearch {
	return &ProwlSearch{p}
}

func (p *ProwlSearch) Go() {
	var fn func(bool)

	fn = func(headless bool) {
		pw, err := NewPW(p.Prowl, &headless)
		if err != nil {
			log.Warn().Err(err).Msg("Prowl - Failed to initialize PW")
			return
		}
		defer pw.Close()

		if err = p.pwWorker(pw); err != nil && headless {
			log.Warn().Err(err).Msg("Prowl - Headless Search Failed; retrying with head ...")
			fn(false)
			return
		}

		if err != nil {
			log.Warn().Err(err).Bool("headless", headless).Msg("Prowl - Search Failed!")
		} else {
			log.Info().Bool("headless", headless).Msg("Prowl - Search Successful!")
		}
	}

	fn(true)
}

func (p *ProwlSearch) pwWorker(pw *PW) (err error) {

	defer pw.Close()

	var page playwright.Page
	var res playwright.Response

	log.Info().Msg("Prowl - Googling ... ")

	if page, err = pw.NewPage(); err != nil {
		return
	} else if res, err = pw.GoTo(page, "https://www.google.com"); err != nil {
		return
	} else if err = pw.IsBlocked(page.URL(), res.URL()); err != nil {
		return
	} else if err = pw.Click(page, googleSearchInputSelectors...); err != nil {
		return
	} else if err = pw.Type(page, pw.Prowler.ID); err != nil {
		return
	} else if err = pw.Press(page, "Enter"); err != nil {
		return
	} else if err = pw.WaitForLoadState(page); err != nil {
		return
	} else if err = pw.IsBlocked(page.URL()); err != nil {
		return
	}

	log.Info().Msg("Prowl - Google Reached!")

	pw.Save(page)

	targetCount := len(pw.Prowler.Targets)
	log.Info().Msgf("PW - Targets [%d]", targetCount)

	if targetCount == 0 {
		return
	}

	var locators []playwright.Locator

	if locators, err = page.Locator(fmt.Sprintf(`[data-dtld]`), playwright.PageLocatorOptions{}).All(); err != nil {
		return
	}

	if len(locators) == 0 {
		log.Warn().Msg("Prowl - No Target Locators Found")
	}

	var att string
	for _, l := range locators {

		att, err = l.GetAttribute("data-dtld")
		if err != nil {
			log.Warn().Err(err).Msg("Prowl - Failed to get Target Locator Attribute")
			continue
		}

		log.Debug().Str("found", att).Msg("Prowl - Locator")
		if !pw.Prowler.Targets.Follow(att) {
			continue
		}

		log.Info().Msgf("Prowl - Target Found [%s]", att)

		if page, err = pw.NewPage(func() error { return l.Click() }); err == nil {
			pw.Save(page)
			err = errors.Join(page.Close())
		}
	}

	return nil
}
