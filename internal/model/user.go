package model

import (
	"bytelyon-functions/internal/service/s3"

	"github.com/oklog/ulid/v2"
)

type User struct {
	ID ulid.ULID `json:"id"`
}

func (u *User) Path() string {
	return "user/" + u.ID.String()
}

func (u *User) Key() string {
	return u.Path() + "/_.json"
}

func FindUser(s string) (*User, error) {

	c, err := NewBasicAuth(s)
	if err != nil {
		return nil, err
	}

	var e *Email
	var p *Password
	if e, err = NewEmail(c.Username); err != nil {
		return nil, err
	} else if p, err = NewPassword(c.Password); err != nil {
		return nil, err
	}

	db := s3.New()
	if err = e.Find(db); err != nil {
		return nil, err
	} else if err = p.Find(db, e.User()); err != nil {
		return nil, err
	}

	return e.User(), nil
}
