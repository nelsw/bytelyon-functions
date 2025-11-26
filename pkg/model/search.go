package model

import (
	"bytelyon-functions/pkg/service/s3"
	"encoding/json"
	"errors"

	"github.com/oklog/ulid/v2"
)

type Search struct {
	*User    `json:"-"`
	ID       ulid.ULID `json:"id"`
	Query    string    `json:"query"`
	VisitAds bool      `json:"visit_ads"`
	Ignore   []string  `json:"ignore"`
	Pages    []Page    `json:"pages"`
}

func (s *Search) Path() string {
	return s.User.Path() + "/search"
}

func (s *Search) Key() string {
	return s.Path() + "/" + s.ID.String() + "/_.json"
}

func (s *Search) WithUser(user *User) *Search {
	s.User = user
	return s
}

func (s *Search) Create() error {

	b, err := json.Marshal(s)
	if err != nil {
		return err
	}

	if s.ID.IsZero() {
		s.ID = NewUlid()
	}

	if err = s3.New().Put(s.Key(), b); err != nil {
		return err
	}

	return nil
}

func (s *Search) FindAll() ([]*Search, error) {

	db := s3.New()

	keys, err := db.Keys(s.Path(), "", 1000)
	if err != nil {
		return nil, err
	}

	var vv []*Search
	for _, k := range keys {

		o, e := db.Get(k)
		if e != nil {
			err = errors.Join(err, e)
			continue
		}

		var v Search
		if e = json.Unmarshal(o, &v); e != nil {
			err = errors.Join(err, e)
			continue
		}

		vv = append(vv, &v)
	}

	return vv, err
}
