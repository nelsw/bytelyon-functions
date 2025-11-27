package model

import (
	"bytelyon-functions/pkg/service/s3"
	"encoding/json"
	"errors"
	"fmt"
)

type Profile struct {
	User *User  `json:"-"`
	Name string `json:"name"`
}

func NewProfile(u *User) *Profile {
	return &Profile{User: u}
}

func (p *Profile) Path() string {
	return p.User.Dir() + "/profile"
}

func (p *Profile) Key() string {
	return p.Path() + "/_.json"
}

func (p *Profile) Validate() error {
	if p.Name == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}

func (p *Profile) Hydrate(a any) error {

	switch t := a.(type) {
	case Profile:
		u := p.User
		*p = t
		p.User = u
	case []byte:
		if err := json.Unmarshal(t, p); err != nil {
			return err
		}
	case string:
		if err := json.Unmarshal([]byte(t), p); err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("unexpected type [%v]", t))
	}
	return nil
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

func (p *Profile) Find(dbs ...s3.Service) error {
	var db s3.Service
	if len(dbs) > 0 {
		db = dbs[0]
	} else {
		db = s3.New()
	}
	return db.Find(p.Key(), p)
}
