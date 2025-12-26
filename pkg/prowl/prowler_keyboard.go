package prowl

import (
	. "bytelyon-functions/pkg/util"

	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
)

func (p *Prowler) Type(s string) error {
	err := p.Page.Keyboard().Type(s, playwright.KeyboardTypeOptions{
		Delay: Ptr(Between(500.0, 1000.0)),
	})
	log.Err(err).Str("text", s).Msg("Prowler - Keyboard#Type")
	return err
}

func (p *Prowler) Press(s string) (err error) {

	err = p.Page.Keyboard().Press(s, playwright.KeyboardPressOptions{
		Delay: Ptr(Between(200, 500.0)),
	})
	log.Err(err).Str("key", s).Msg("Prowler - Keyboard#Press")
	if err != nil {
		return err
	}

	err = p.Page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State:   playwright.LoadStateNetworkidle,
		Timeout: Ptr(60_000.0),
	})
	log.Err(err).Msg("Prowler - Wait for load")
	if err != nil {
		return err
	}

	log.Info().Bool("headless", p.headless).Str("s", s).Msg("Press")
	return nil
}
