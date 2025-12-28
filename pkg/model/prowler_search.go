package model

import (
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func (p *Prowler) ProwlSearch() ulid.ULID {

	var fn func(ulid.ULID, bool)

	fn = func(prowlID ulid.ULID, headless bool) {
		pw, err := NewPW(p, &headless)
		if err != nil {
			log.Warn().Err(err).Msg("Prowler - PW failed to initialize")
			return
		}
		defer pw.Close()

		log.Info().Msg("Prowler - Searching ... ")
		if err = pw.Search(prowlID); err != nil && headless {
			log.Warn().Err(err).Msg("Prowler - Headless Search Failed; retrying with head ...")
			fn(prowlID, false)
			return
		}

		if err != nil {
			log.Warn().Err(err).Msg("Prowler - Headed Search Failed!")
		} else {
			log.Info().Bool("headless", headless).Msg("Prowler - Search Succeeded")
		}
	}

	prowlID := NewUlid()

	fn(prowlID, true)

	return prowlID
}

func handleProwlSearch(prowlID ulid.ULID, headless bool) {

}
