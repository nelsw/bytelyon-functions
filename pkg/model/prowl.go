package model

import (
	"bytelyon-functions/pkg/db"
	"bytelyon-functions/pkg/util"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type Prowl struct {
	Prowler *Prowler  `json:"prowler"`
	ID      ulid.ULID `json:"id"`
}

func NewProwl(p *Prowler) *Prowl {
	return &Prowl{p, NewUlid()}
}

func (p *Prowl) String() string {
	return p.Prowler.String() + util.Path("prowl", p.ID)
}

func (p *Prowl) Go() {

	switch p.Prowler.Type {
	case SearchProwlerType:
		NewProwlSearch(p).Go()
	case SitemapProwlerType:
		NewProwlSitemap(p).Go()
	case NewsProwlerType:
		NewProwlNews(p).Go()
	}

	p.Prowler.Prowled = p.ID
	p.Prowler.Duration = time.Since(p.ID.Timestamp())
	if err := db.Save(p.Prowler); err != nil {
		log.Warn().Err(err).Msgf("Prowl - Failed to save Prowler [%s]", p)
	}
}
