package model

import (
	"bytelyon-functions/pkg/service/em"
	"bytelyon-functions/pkg/service/s3"
	"encoding/json"
	"regexp"
	"strings"

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

func (u *User) String() string {
	return "user/" + u.ID.String()
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

func (u *User) Searches() ([]*Search, error) {

	//search := Search{User: u}
	//r := regexp.MustCompile(`.*user/[A-Za-z0-9]{26}/search/[A-Za-z0-9]{26}/result/[A-Za-z0-9]{26}/page/[A-Za-z0-9]{26}/_.json`)
	//searches, err := em.FindAll(&search, r)
	//if err != nil || len(searches) == 0 {
	//	return searches, err
	//}
	//log.Debug().EmbedObject(searches[0]).Msg("found search")
	//log.Err(err).
	//	Int("searches", len(searches)).
	//	Msg("find searches")
	//
	//for _, s := range searches {
	//	if err = s.FetchPages(u); err != nil {
	//		return searches, err
	//	}
	//}

	db := s3.New()
	arr, err := db.Keys(`user/`+u.ID.String()+`/search/`, "", 1000)
	if err != nil {
		panic(err)
	}

	//sr := regexp.MustCompile(`/search/[A-Za-z0-9]{26}/result/[A-Za-z0-9]{26}/page/[A-Za-z0-9]{26}/_.json`)
	sr := regexp.MustCompile(`.*/search/[A-Za-z0-9]{26}/_.json`)
	var searchKeys []string
	for _, k := range arr {
		if sr.MatchString(k) {
			searchKeys = append(searchKeys, k)
		}
	}

	findPages := func(k string) []*Page {
		k = strings.TrimSuffix(k, "_.json")
		//k += "/result/"
		// regexp.MustCompile(`/search/[A-Za-z0-9]{26}/result/[A-Za-z0-9]{26}/page/[A-Za-z0-9]{26}/_.json`)
		keys, err := db.Keys(k, "", 1000)
		if err != nil {
			panic(err)
		}
		pr := regexp.MustCompile(`/search/[A-Za-z0-9]{26}/result/[A-Za-z0-9]{26}/page/[A-Za-z0-9]{26}/_.json`)
		var pages []*Page
		for _, k = range keys {
			if !pr.MatchString(k) {
				continue
			}
			var p Page
			b, err := db.Get(k)
			err = json.Unmarshal(b, &p)
			if err != nil {
				panic(err)
			}
			url := strings.TrimSuffix(k, "/_.json")
			if out, err := db.GetPresigned(url + "/content.html"); err == nil {
				p.Content = out
			}

			if out, err := db.GetPresigned(url + "/screenshot.png"); err == nil {
				p.Screenshot = out
			}

			if out, err := db.Get(url + "/results.json"); err == nil {
				if err = json.Unmarshal(out, &p.Results); err != nil {
					log.Warn().Err(err).Msg("failed to unmarshal pagee results")
				}
			}
			pages = append(pages, &p)
		}
		return pages

	}
	log.Trace().Msgf("searches found: %d", len(searchKeys))
	var searches []*Search
	for _, k := range searchKeys {
		var s Search
		b, err := db.Get(k)
		err = json.Unmarshal(b, &s)
		if err != nil {
			panic(err)
		}
		s.UserID = u.ID
		s.Pages = findPages(k)
		log.Trace().EmbedObject(&s).Msg("found search")
		searches = append(searches, &s)
	}
	//
	//r := regexp.MustCompile(`/page/[A-Za-z0-9]{26}/_.json`)
	//var keys []string
	//for _, k := range arr {
	//	if r.MatchString(k) {
	//		fmt.Println(k)
	//		keys = append(keys, k)
	//	}
	//}
	//
	//for _, k := range keys {
	//
	//}

	return searches, err
}

func (u *User) Sitemaps() ([]*Sitemaps, error) {
	sitemap := Sitemap{User: u}
	sitemaps, err := em.FindAll(&sitemap, regexp.MustCompile(sitemap.Path()+`/[A-Za-z0-9\\.]+/[A-Za-z0-9]{26}/_.json`))
	log.Err(err).Int("sitemaps", len(sitemaps)).Msg("find sitemaps")
	if err != nil {
		return nil, err
	}

	m := make(map[string][]*Sitemap)
	for _, s := range sitemaps {
		m[s.Domain] = append(m[s.Domain], s)
	}

	var result []*Sitemaps
	for _, v := range m {
		result = append(result, NewSitemaps(v))
	}
	return result, nil
}
