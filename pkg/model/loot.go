package model

import (
	"bytelyon-functions/pkg/service/s3"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Loot struct {
	*Plunder `json:"-"`
	IDX      ulid.ULID `json:"idx"`
	ID       ulid.ULID `json:"id"`
	Name     string    `json:"name"`
	Data     any       `json:"data"`
	HTML     string    `json:"html"`
	Image    string    `json:"image"`
}

func (l *Loot) MarshalZerologObject(evt *zerolog.Event) {
	evt.Str("id", l.ID.String()).
		Str("plunder", l.Plunder.ID.String()).
		Str("name", l.Name).
		Bool("data", l.Data != nil).
		Bool("html", l.HTML != "").
		Bool("image", l.Image != "")
}

func NewLoot(p *Plunder, id, ts string, titles []string) *Loot {

	urlFn := func(s string) string {
		if u, err := s3.New().GetPresigned(s); err != nil {
			log.Err(err).Msg("failed to get presigned url")
			return s
		} else {
			return u
		}
	}

	dataFn := func(s string) map[string]any {
		var m map[string]any
		if b, err := s3.New().Get(s); err != nil {
			log.Err(err).Msg("failed to get loot")
		} else {
			_ = json.Unmarshal(b, &m)
		}
		return m
	}

	var name string
	var data any
	var html, image string
	for _, s := range titles {

		if name == "" {
			name = strings.Split(s, ".")[0]
		}

		key := p.Dir() + "/loot/" + id + "/" + ts + " " + s
		if strings.HasSuffix(s, ".html") {
			html = urlFn(key)
		} else if strings.HasSuffix(s, ".png") {
			image = urlFn(key)
		} else if strings.HasSuffix(s, ".json") {
			data = dataFn(key)
		}
	}
	var unix int64
	if TS, err := strconv.Atoi(ts); err != nil {
		log.Err(err).Msg("failed to parse timestamp!?")
		unix = time.Now().UTC().UnixMilli()
	} else {
		unix = int64(TS)
	}

	return &Loot{
		Plunder: p,
		IDX:     ulid.MustParse(id),
		ID:      NewUlidFromTime(time.Unix(unix, 0)),
		Name:    name,
		Data:    data,
		HTML:    html,
		Image:   image,
	}
}
