package model

import (
	"bytelyon-functions/pkg/service/em"
	"bytelyon-functions/pkg/service/s3"
	"encoding/json"
	"errors"

	"github.com/rs/zerolog/log"
)

type Profile struct {
	User *User  `json:"-"`
	Name string `json:"name"`
}

func (p *Profile) Path() string {
	return p.User.Dir() + "/profile"
}

func (p *Profile) Dir() string {
	return p.Path()
}

func (p *Profile) Key() string {
	return p.Dir() + "/_.json"
}

func (p *Profile) Create(b []byte) (*Profile, error) {
	var v Profile
	if err := json.Unmarshal(b, &v); err != nil {
		log.Err(err).Msg("failed to unmarshal profile")
		return nil, err
	}
	if v.Name == "" {
		return nil, errors.New("name cannot be empty")
	}
	v.User = p.User
	if err := em.Save(&v); err != nil {
		log.Err(err).Msg("failed to save profile")
		return nil, err
	}
	return &v, nil
}

func (p *Profile) Save(dbs ...s3.Service) error {

	var db s3.Service
	if len(dbs) > 0 {
		db = dbs[0]
	} else {
		db = s3.New()
	}

	if b, err := json.Marshal(p); err != nil {
		return err
	} else if err = db.Put(p.Key(), b); err != nil {
		return err
	}

	return nil
}

func (p *Profile) Find() (*Profile, error) {
	user := p.User
	if err := em.Find(p); err != nil {
		return nil, err
	}
	p.User = user
	return p, nil
}

func NewProfile(u *User) *Profile {
	return &Profile{User: u}
}
