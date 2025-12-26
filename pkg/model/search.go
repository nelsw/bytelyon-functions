package model

import (
	"bytelyon-functions/pkg/service/em"
	"bytelyon-functions/pkg/service/s3"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Search struct {
	*User   `json:"-"`
	UserID  ulid.ULID `json:"user_id"`
	ID      ulid.ULID `json:"id"`
	Query   string    `json:"query"`
	Targets Targets   `json:"targets"`
	Pages   []*Page   `json:"pages,omitempty"`
}

func (s *Search) Path() string {
	return s.User.Dir() + "/search"
}

func (s *Search) Dir() string {
	return s.Path() + "/" + s.ID.String()
}

func (s *Search) Key() string {
	return s.Dir() + "/_.json"
}

func (s *Search) MarshalZerologObject(evt *zerolog.Event) {
	if s.User != nil {
		evt.EmbedObject(s.User)
	}
	evt.Str("search", s.ID.String()[20:]).
		Str("user", s.UserID.String()[20:]).
		Str("query", s.Query).
		Int("pages", len(s.Pages)).
		EmbedObject(s.Targets)
}

func (s *Search) HasTargets() bool {
	return s.Targets != nil && !s.Targets.None()
}

func (s *Search) IsTarget(t string) bool {
	return s.HasTargets() && s.Targets.Exist(t)
}

func (s *Search) CreatePage(ts time.Time, u, t, c string, ss []byte) {

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
	return s3.New().Put(s, a, "/state.json")
}

func (s *Search) FindState() *playwright.OptionalStorageState {

	var state playwright.OptionalStorageState
	if out, err := s3.New().Get(s.String() + "/state.json"); err == nil {
		err = json.Unmarshal(out, &state)
	}

	if len(state.Cookies) > 100 {
		state = playwright.OptionalStorageState{}
	}

	return &state
}

func (s *Search) Create(b []byte) (*Search, error) {

	var v Search
	if err := json.Unmarshal(b, &v); err != nil {
		log.Err(err).Msg("failed to unmarshal search")
		return nil, err
	}

	if v.Query == "" {
		return nil, errors.New("missing query")
	}

	v.User = s.User
	if v.ID.IsZero() {
		v.ID = NewUlid()
	}

	err := em.Save(&v)
	log.Err(err).EmbedObject(&v).Msg("save search")

	return &v, err
}

func (s *Search) Delete() error {

	if err := em.Delete(s); err != nil {
		return err
	}

	// delete the associated job
	j, err := NewJob(s.User, s.ID).Find()
	if err != nil {
		log.Warn().Err(err).Msg("failed to find search, it may not exist or have been deleted")
		return nil
	}

	return em.Delete(j)
}

func NewSearch(args ...any) *Search {
	var v = new(Search)
	for _, arg := range args {
		switch arg.(type) {
		case *User:
			v.User = arg.(*User)
		case ulid.ULID:
			v.ID = arg.(ulid.ULID)
		case string:
			if id, err := ulid.Parse(arg.(string)); err == nil {
				v.ID = id
			}
		}
	}
	return v
}
