package model

import (
	"bytelyon-functions/pkg/service/em"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Targets map[string]bool

func (t Targets) MarshalZerologObject(evt *zerolog.Event) {
	var s []string
	for k, v := range t {
		s = append(s, fmt.Sprintf("%s:%v", k, v))
	}
	evt.Strs("targets", s)
}

type Search struct {
	*User   `json:"-"`
	ID      ulid.ULID `json:"id"`
	Query   string    `json:"query"`
	Targets Targets   `json:"targets"`
	Pages   []*Page   `json:"pages"`
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
	evt.Stringer("id", s.ID).
		Str("query", s.Query).
		EmbedObject(s.Targets)
}

func (s *Search) FetchPages(u *User) (err error) {

	s.User = u

	page := NewPage(s)
	s.Pages, err = em.FindAll(page, regexp.MustCompile(page.Path()+`/[A-Za-z0-9]{26}/_.json`))

	for _, p := range s.Pages {
		p.FetchData(s)
	}

	log.Err(err).
		Int("pages", len(s.Pages)).
		Msg("fetch pages")

	return err
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
