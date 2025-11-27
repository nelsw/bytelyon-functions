package model

import (
	"bytelyon-functions/pkg/service/em"
	"regexp"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type Article struct {
	*News  `json:"-"`
	ID     ulid.ULID `json:"id"`
	URL    string    `json:"url"`
	Title  string    `json:"title"`
	Time   int64     `json:"time"`
	Source string    `json:"source"`
}

func (a *Article) Path() string {
	return a.News.Dir() + "/article"
}

func (a *Article) Dir() string {
	return a.Path() + "/" + a.ID.String()
}

func (a *Article) Key() string {
	return a.Dir() + "/_.json"
}

func (a *Article) Delete() error {
	return em.Delete(a)
}

func (a *Article) Find() error {
	return em.Find(a)
}

func (a *Article) FindAll() ([]*Article, error) {
	all, err := em.FindAll(a, regexp.MustCompile(`.*/article/([A-Za-z0-9]{26}/_.json)$`))
	if err != nil {
		log.Err(err).Msg("failed to find articles")
		return []*Article{}, err
	}
	for i, v := range all {
		v.News = a.News
		all[i] = v
	}
	return all, nil
}

func (a *Article) FindLast() error {
	return em.FindLast(a)
}

func (a *Article) IDs() ([]string, error) {
	keys, err := em.Keys(a, regexp.MustCompile(`^`+a.Path()+`/([A-Za-z0-9]{26}/_.json)$`))
	if err != nil {
		return nil, err
	}
	for i, k := range keys {
		keys[i] = k[len(a.Path())+1 : len(k)-len("/_.json")]
	}
	return keys, nil
}

func NewArticle(u *User, ids ...string) *Article {

	var a = new(Article)
	if len(ids) < 1 {
		return a
	}

	if a.News = NewNews(u, ids[0]); len(ids) == 1 {
		return a
	}

	if id, err := ulid.Parse(ids[1]); err == nil {
		a.ID = id
	}

	return a
}
