package model

import (
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func (p *Prowler) ProwlSearch(headless bool) ulid.ULID {

	prowl, err := NewProwl(p, &headless)
	if err != nil {
		log.Warn().Err(err).Msg("Prowler - Prowl failed to initialize")
		return prowl.ID
	}
	defer prowl.Close()

	log.Info().Msg("Prowler - Searching ... ")

	if err = prowl.Search(); err != nil && headless {
		log.Warn().Err(err).Msg("Prowler - Headless Search Failed; retrying with head ...")
		return p.ProwlSearch(false)
	}

	if err != nil {
		log.Warn().Err(err).Msg("Prowler - Headed Search Failed!")
	} else {
		log.Info().Bool("headless", headless).Msg("Prowler - Search Succeeded")
	}

	return prowl.ID
}
