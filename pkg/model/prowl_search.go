package model

import (
	"errors"
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

type ProwlSearch struct {
	*Prowl
}

func NewProwlSearch(p *Prowl) *ProwlSearch {
	return &ProwlSearch{p}
}

func (p *ProwlSearch) Go() ulid.ULID {
	var fn func(bool) ulid.ULID

	fn = func(headless bool) ulid.ULID {

		pw, err := NewPW(p.Prowl, &headless)
		if err != nil {
			log.Warn().Err(err).Msg("ProwlSearch - Failed to initialize PW")
			return ulid.Zero
		}
		defer pw.Close()

		log.Info().Bool("headless", headless).Msg("ProwlSearch - Working ... ")

		var prowled ulid.ULID
		if prowled, err = p.pwWorker(pw); err != nil && headless {
			log.Warn().Err(err).Msg("ProwlSearch - Headless Failed; retrying with head ...")
			return fn(false)
		}

		if err != nil {
			log.Warn().Err(err).Bool("headless", headless).Msg("ProwlSearch - Failed!")
			return ulid.Zero
		}

		log.Info().Bool("headless", headless).Msg("ProwlSearch - Success!")
		return prowled
	}

	return fn(true)
}

func (p *ProwlSearch) pwWorker(pw *PW) (prowled ulid.ULID, err error) {

	defer pw.Close()

	var page playwright.Page
	var res playwright.Response

	if page, err = pw.NewPage(); err != nil {
		return
	} else if res, err = pw.GoTo(page, "https://www.google.com"); err != nil {
		return
	} else if err = pw.IsBlocked(page, res); err != nil {
		return
	} else if err = pw.Click(page, googleSearchInputSelectors...); err != nil {
		return
	} else if err = pw.Type(page, pw.Prowler.ID); err != nil {
		return
	} else if err = pw.Press(page, "Enter"); err != nil {
		return
	} else if err = pw.WaitForLoadState(page); err != nil {
		return
	} else if err = pw.IsBlocked(page); err != nil {
		return
	}

	prowled = pw.Save(page)

	targetCount := len(pw.Prowler.Targets)
	log.Info().Msgf("ProwlSearch - Targets [%d]", targetCount)

	if targetCount == 0 {
		return
	}

	var locators []playwright.Locator
	if locators, err = page.Locator(fmt.Sprintf(`[data-dtld]`), playwright.PageLocatorOptions{}).All(); err != nil {
		return
	}

	if len(locators) == 0 {
		log.Warn().Msg("ProwlSearch - No Target Locators Found")
		return
	}

	var att string
	for _, l := range locators {

		if att, err = l.GetAttribute("data-dtld"); err != nil {
			log.Warn().Err(err).Msg("ProwlSearch - Failed to get Target Locator Attribute")
			continue
		}

		log.Debug().Str("found", att).Msg("ProwlSearch - Locator")
		if !pw.Prowler.Targets.Follow(att) {
			continue
		}

		log.Info().Msgf("ProwlSearch - Target Found [%s]", att)
		if page, err = pw.NewPage(func() error { return l.Click() }); err == nil {
			prowled = pw.Save(page)
			err = errors.Join(page.Close())
		}
	}

	return prowled, nil
}
