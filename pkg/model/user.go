package model

import (
	"bytelyon-functions/pkg/service/em"
	"bytelyon-functions/pkg/service/s3"
	"regexp"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var (
	userKeyRegex = regexp.MustCompile(`.*user/([A-Za-z0-9]{26}/_.json)$`)
)

type User struct {
	ID ulid.ULID `json:"id"`
}

func MakeDemoUser() User {
	return User{ID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0")}
}

func (u *User) Path() string {
	return "user/"
}

func (u *User) Dir() string {
	return u.Path() + u.ID.String()
}

func (u *User) Key() string {
	return u.Dir() + "/_.json"
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

func FindAllUsers() ([]*User, error) {
	users, err := em.FindAll(&User{}, userKeyRegex)
	if err != nil {
		log.Err(err).Msg("failed to find users")
		return []*User{}, err
	}
	return users, nil
}
