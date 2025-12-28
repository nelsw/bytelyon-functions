package model

import (
	"github.com/rs/zerolog/log"
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
			log.Warn().Err(err).Msg("Prowler - PW failed to initialize")
			return
		}
		defer pw.Close()

		log.Info().Msg("Prowler - Searching ... ")
		if err = pw.Search(); err != nil && headless {
			log.Warn().Err(err).Msg("Prowler - Headless Search Failed; retrying with head ...")
			fn(false)
			return
		}

		if err != nil {
			log.Warn().Err(err).Msg("Prowler - Headed Search Failed!")
		} else {
			log.Info().Bool("headless", headless).Msg("Prowler - Search Succeeded")
		}
	}
}
