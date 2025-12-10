package model

import (
	"math"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
)

type Sitemaps struct {
	ID       ulid.ULID  `json:"id"`
	URL      string     `json:"url"`
	Domain   string     `json:"domain"`
	Updated  int64      `json:"updated"`
	Duration float64    `json:"duration"`
	Average  float64    `json:"average"`
	Relative int        `json:"relative"`
	Remote   int        `json:"remote"`
	Size     int        `json:"size"`
	Sitemaps []*Sitemap `json:"sitemaps"`
}

func (s *Sitemaps) MarshalZerologObject(evt *zerolog.Event) {
	evt.Stringer("sitemaps", s.ID).
		Time("created", s.ID.Timestamp().UTC()).
		Time("updated", time.UnixMilli(s.Updated).UTC()).
		Int("size", s.Size).
		Str("url", s.URL).
		Str("domain", s.Domain).
		Float64("duration", s.Duration).
		Int("relative", s.Relative).
		Int("remote", s.Remote).
		Float64("average", s.Average)
}

func NewSitemaps(sitemaps []*Sitemap) *Sitemaps {

	s := Sitemaps{Size: len(sitemaps)}

	for i, sitemap := range sitemaps {
		s.Sitemaps = append(s.Sitemaps, sitemap)

		if i == 0 {
			s.ID = sitemap.ID
			s.URL = sitemap.URL
			s.Domain = sitemap.Domain
		}

		if i == len(sitemaps)-1 {
			s.Updated = sitemap.ID.Timestamp().UnixMilli()
		}

		s.Duration += sitemap.Duration
		s.Relative += len(sitemap.Relative)
		s.Remote += len(sitemap.Remote)
	}

	s.Average = math.Trunc((s.Duration/float64(len(sitemaps)))*100) / 100

	return &s
}
