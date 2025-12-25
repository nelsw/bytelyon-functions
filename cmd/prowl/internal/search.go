package internal

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Search struct {
	UserID   ulid.ULID `json:"user_id"`
	ResultID ulid.ULID `json:"result_id"`
	ID       ulid.ULID `json:"id"`
	Query    string    `json:"query"`
	Targets  *Targets  `json:"targets"`
}

func (s *Search) String() string {
	return Path("user", s.UserID, "search", s.ID)
}

func (s *Search) MarshalZerologObject(evt *zerolog.Event) {
	evt.Stringer("user", s.UserID).
		Stringer("search", s.ID).
		Str("query", s.Query)
}

func (s *Search) HasTargets() bool {
	return s.Targets != nil && !s.Targets.None()
}

func (s *Search) IsTarget(t string) bool {
	return s.HasTargets() && s.Targets.Exist(t)
}

func (s *Search) CreatePage(ts time.Time, u, t, c string, ss []byte) {
	log.Trace().EmbedObject(s).Msg("create page")
	p := &Page{s.UserID, s.ID, NewULID(ts), NewULID(), u, t}
	_ = S3.Put(p)
	_ = S3.Put(p, ss, "/screenshot.png")
	_ = S3.Put(p, c, "/content.html")
	if !strings.HasPrefix(p.URL, "https://www.google.com") {
		return
	}
	if results, err := NewResults(c); err == nil {
		_ = S3.Put(p, results.Data, "/results.json")
	}
	log.Info().EmbedObject(p).Msg("created page")
}

func (s *Search) SaveState(a any, err error) error {
	if err != nil {
		log.Warn().
			Err(err).
			Msg("failed to read prowl browser context storage state")
	}
	if a == nil {
		return nil
	}
	return S3.Put(s, a, "/state.json")
}

func (s *Search) FindState() *playwright.OptionalStorageState {

	var state playwright.OptionalStorageState
	if out, err := S3.Get(s.String() + "/state.json"); err == nil {
		err = json.Unmarshal(out, &state)
	}

	if len(state.Cookies) > 100 {
		state = playwright.OptionalStorageState{}
	}

	return &state
}
