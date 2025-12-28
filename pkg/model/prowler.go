package model

import (
	"bytelyon-functions/pkg/util"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type Prowler struct {
	UserID ulid.ULID `json:"user_id"`
	// ID is either a URL or Query string
	ID        string        `json:"id"`
	Type      ProwlerType   `json:"type"`
	Prowled   ulid.ULID     `json:"prowled"`
	Frequency time.Duration `json:"frequency"`
	Duration  time.Duration `json:"duration"`
	Targets   Targets       `json:"targets,omitempty"`
}

func (p *Prowler) String() string {
	id := p.ID
	if p.Type == SitemapProwlerType {
		id = util.Domain(p.ID)
	}
	return fmt.Sprintf("user/%s/prowler/%s/%s", p.UserID, p.Type, id)
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

	NewProwl(p).Go()
}
