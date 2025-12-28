package model

import (
	"bytelyon-functions/pkg/service/s3"
	"encoding/json"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type Prowler struct {
	UserID  ulid.ULID   `json:"user_id"`
	ID      ulid.ULID   `json:"id"`
	Type    ProwlerType `json:"type"`
	Targets Targets     `json:"targets"`
	Query   string      `json:"query"`
	Prowled ulid.ULID   `json:"prowled"`
}

func NewProwler(a ...any) *Prowler {
	var p = new(Prowler)
	for _, v := range a {
		switch v.(type) {
		case ulid.ULID:
			if p.UserID.IsZero() {
				p.UserID = v.(ulid.ULID)
			} else {
				p.ID = v.(ulid.ULID)
			}
		case Targets:
			p.Targets = v.(Targets)
		case ProwlerType:
			p.Type = v.(ProwlerType)
		case string:
			if id, err := ulid.Parse(v.(string)); err == nil && p.ID.IsZero() {
				p.ID = id
				continue
			} else if p.Type == "" {
				p.Type = v.(ProwlerType)
			} else {
				p.Query = v.(string)
			}
		}
	}

	if p.ID.IsZero() {
		p.ID = NewUlid()
	}

	return p
}

func (p *Prowler) Prowl(a ...any) {

	headless := len(a) > 0 && a[0].(bool)

	prowl, err := NewProwl(p, &headless)
	if err != nil {
		log.Warn().Err(err).Msg("Prowler - Prowl failed to initialize")
		return
	}
	defer prowl.Close()

	if p.Type == SearchProwlType {

		log.Info().Msg("Prowler - Searching ... ")

		if err = prowl.Search(); err != nil && headless {
			log.Warn().Err(err).Msg("Prowler - Headless Search Failed; retrying with head ...")
			p.Prowl()
			return
		}

		if err != nil {
			log.Warn().Err(err).Msg("Prowler - Headed Search Failed!")
		} else {
			log.Info().Bool("headless", headless).Msg("Prowler - Search Succeeded")
		}
	}

	p.Prowled = prowl.ID
	b, _ := json.Marshal(p)
	s3.New().Put("user/"+p.UserID.String()+"/prowler/search/"+p.ID.String()+"/_.json", b)
}
