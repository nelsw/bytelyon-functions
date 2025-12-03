package model

import (
	"bytelyon-functions/pkg/service/s3"
	"strconv"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LootType string

const (
	LootTypePng  LootType = "png"
	LootTypeHtml          = "html"
	LootTypeJson          = "json"
)

type Loot struct {
	*Plunder `json:"-"`
	ID       ulid.ULID `json:"id"`
	Key      string    `json:"-"`
	URL      string    `json:"url"`
	Type     LootType  `json:"type"`
	Name     string    `json:"name"`
}

func (l *Loot) MarshalZerologObject(evt *zerolog.Event) {
	evt.Str("id", l.ID.String()).
		Str("plunder", l.Plunder.ID.String()).
		Str("url", l.URL).
		Str("type", string(l.Type)).
		Str("name", l.Name)
}

func NewLoot(p *Plunder, k string) *Loot {

	lt := LootTypePng
	if strings.HasSuffix(k, ".html") {
		lt = LootTypeHtml
	} else if strings.HasSuffix(k, ".json") {
		lt = LootTypeJson
	}

	lastSlash := strings.LastIndex(k, "/")
	left, right, _ := strings.Cut(k[lastSlash+1:], " ")

	var id ulid.ULID
	if i, err := strconv.Atoi(left); err == nil {
		id = NewUlidFromTime(time.Unix(int64(i), 0))
	}

	url, err := s3.New().GetPresigned(k)
	if err != nil {
		log.Err(err).Msg("failed to get presigned url")
		url = k
	}

	return &Loot{
		Plunder: p,
		ID:      id,
		Key:     k,
		URL:     url,
		Name:    strings.TrimSuffix(right, `.`+string(lt)),
		Type:    lt,
	}
}
