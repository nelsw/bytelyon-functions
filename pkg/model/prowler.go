package model

import (
	"bytelyon-functions/pkg/db"
	"bytelyon-functions/pkg/util"
	"encoding/json"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type Prowler struct {
	UserID    ulid.ULID     `json:"user_id"`
	ID        ulid.ULID     `json:"id"`
	Type      ProwlerType   `json:"type"`
	Prowled   ulid.ULID     `json:"prowled"`
	Frequency time.Duration `json:"frequency"`
	Duration  time.Duration `json:"duration"`
	Targets   Targets       `json:"targets,omitempty"`
	Query     string        `json:"query,omitempty"`
	URL       string        `json:"url,omitempty"`
	Domain    string        `json:"domain,omitempty"`
	Relative  []string      `json:"relative,omitempty"`
	Remote    []string      `json:"remote,omitempty"`
}

func (p *Prowler) String() string {
	return util.Path("user", p.UserID, "prowler", p.Type, "prowl", p.ID)
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

	var id ulid.ULID
	switch p.Type {
	case SearchProwlerType:
		id = p.ProwlSearch()
	case SitemapProwlerType:
		id = p.ProwlSitemap()
	case NewsProwlerType:
		id = p.ProwlNews()
	}

	p.Prowled = id
	p.Duration = time.Since(id.Timestamp())
	if err := db.Save(p); err != nil {
		b, _ := json.MarshalIndent(p, "", "\t")
		log.Warn().Err(err).Msg("Prowler - Failed to save prowler: " + string(b))
	}
}
