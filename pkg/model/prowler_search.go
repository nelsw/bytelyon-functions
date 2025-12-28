package model

import (
	"bytelyon-functions/pkg/db"

	"github.com/rs/zerolog/log"
)

func (p *Prowler) Search(headless bool) {

	prowl, err := NewProwl(p, &headless)
	if err != nil {
		log.Warn().Err(err).Msg("Prowler - Prowl failed to initialize")
		return
	}
	defer prowl.Close()

	log.Info().Msg("Prowler - Searching ... ")

	if err = prowl.Search(); err != nil && headless {
		log.Warn().Err(err).Msg("Prowler - Headless Search Failed; retrying with head ...")
		p.Search(false)
		return
	}

	if err != nil {
		log.Warn().Err(err).Msg("Prowler - Headed Search Failed!")
	} else {
		log.Info().Bool("headless", headless).Msg("Prowler - Search Succeeded")
	}

	p.Prowled = prowl.ID
	db.Save(p)
}
