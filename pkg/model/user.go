package model

import (
	"bytelyon-functions/pkg/service/em"
	"bytelyon-functions/pkg/service/s3"
	"regexp"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
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

func NewDemoUser() *User {
	u := MakeDemoUser()
	return &u
}

func (u *User) MarshalZerologObject(evt *zerolog.Event) {
	evt.Stringer("user", u.ID)
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

func (u *User) Searches() ([]*Search, error) {

	search := Search{User: u}
	searches, err := em.FindAll(&search, regexp.MustCompile(search.Path()+`/[A-Za-z0-9]{26}/_.json`))

	log.Err(err).
		Int("searches", len(searches)).
		Msg("find searches")

	for _, s := range searches {
		if err = s.FetchPages(u); err != nil {
			return searches, err
		}
	}

	return searches, err
}

func (u *User) Sitemaps() (*Sitemaps, error) {
	sitemap := Sitemap{User: u}
	sitemaps, err := em.FindAll(&sitemap, regexp.MustCompile(sitemap.Path()+`/[A-Za-z0-9\\.]+/[A-Za-z0-9]{26}/_.json`))
	log.Err(err).Int("sitemaps", len(sitemaps)).Msg("find sitemaps")
	if err != nil {
		return nil, err
	}

	return NewSitemaps(sitemaps), nil
}
