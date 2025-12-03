package model

import (
	"bytelyon-functions/pkg/service/s3"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LootType string

const (
	LootTypePng  LootType = "png"
	LootTypeHtml          = "html"
)

type Loot struct {
	Key     string   `json:"key"`
	URL     string   `json:"url"`
	Type    LootType `json:"type"`
	Time    int64    `json:"time"`
	Title   string   `json:"title"`
	Content string   `json:"content,omitempty"`
}

func (l *Loot) MarshalZerologObject(evt *zerolog.Event) {
	evt.Str("key", l.Key).
		Str("url", l.URL).
		Str("type", string(l.Type)).
		Int64("time", l.Time).
		Str("title", l.Title)
}

func NewLoot(k string) *Loot {

	lt := LootTypePng
	if strings.HasSuffix(k, ".html") {
		lt = LootTypeHtml
	}

	var title string
	var time int64
	if idx := strings.Index(k, " "); idx > -1 {
		title = k[idx+1 : len(k)-4]
		if i, err := strconv.Atoi(k[:idx]); err == nil {
			time = int64(i)
		}
	} else {
		title = k
	}

	url, err := s3.New().GetPresigned(k)
	if err != nil {
		log.Err(err).Msg("failed to get presigned url")
		url = k
	}

	var content string
	if lt == LootTypeHtml {
		b, _ := s3.New().Get(k)
		content = string(b)
	}

	return &Loot{
		Key:     k,
		URL:     url,
		Title:   title,
		Type:    lt,
		Time:    time,
		Content: content,
	}
}
