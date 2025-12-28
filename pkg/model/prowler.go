package model

import (
	"bytelyon-functions/pkg/util"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type Prowler struct {
	UserID    ulid.ULID     `json:"user_id"`
	ID        ulid.ULID     `json:"id"`
	Type      ProwlerType   `json:"type"`
	Targets   Targets       `json:"targets,omitempty"`
	Query     string        `json:"query,omitempty"`
	URL       string        `json:"url,omitempty"`
	Frequency time.Duration `json:"frequency"`
	Prowled   ulid.ULID     `json:"prowled"`
}

func (p *Prowler) String() string {
	return util.Path("user", p.UserID, "prowler", p.Type, p.ID)
}

func (p *Prowler) Prowl() {

	if p.Frequency == 0 && !p.Prowled.IsZero() {
		log.Info().Msg("Prowler - Already prowled ...")
		return
	}

	if p.Prowled.Timestamp().Add(p.Frequency).After(time.Now()) {
		log.Info().Msg("Prowler - Too soon to prowl ...")
		return
	}

	switch p.Type {
	case SearchProwlType:
		p.Search(true)
	}
}
