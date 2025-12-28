package model

import (
	"bytelyon-functions/pkg/db"
	"encoding/json"
	"fmt"
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

func (p *Prowl) Go() {

	switch p.Prowler.Type {
	case SearchProwlerType:
		p.Prowler.Prowled = NewProwlSearch(p).Go()
	case SitemapProwlerType:
		p.Prowler.Prowled = NewProwlSitemap(p).Go()
	case NewsProwlerType:
		p.Prowler.Prowled = NewProwlNews(p).Go()
	}

	p.Prowler.Duration = time.Since(p.ID.Timestamp())
	if err := db.Save(p.Prowler); err != nil {
		log.Warn().Err(err).Msgf("Prowl - Failed to save Prowler [%s]", p)
		return
	}

	if p.Prowler.Type == SearchProwlerType {
		b, _ := json.Marshal(p)
		db.NewS3().Put(fmt.Sprintf("%s/%s/_.json", p.Prowler, p.ID), b)
	}
}
