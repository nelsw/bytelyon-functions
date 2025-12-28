package model

import (
	. "bytelyon-functions/pkg/util"

	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
)

func (pw *PW) Type(page playwright.Page, s string) error {
	err := page.Keyboard().Type(s, playwright.KeyboardTypeOptions{
		Delay: Ptr(Between(500.0, 1000.0)),
	})
	log.Err(err).Str("text", s).Msg("PW - Keyboard#Type")
	return err
}

func (pw *PW) Press(page playwright.Page, s string) (err error) {
	err = page.Keyboard().Press(s, playwright.KeyboardPressOptions{
		Delay: Ptr(Between(200, 500.0)),
	})
	log.Err(err).Str("key", s).Msg("PW - Keyboard#Press")
	return err
}
