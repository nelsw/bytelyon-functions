package model

import (
	"bytelyon-functions/pkg/service/em"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type News struct {
	*User    `json:"-"`
	ID       ulid.ULID  `json:"id"`
	Keywords []string   `json:"keywords"`
	Articles []*Article `json:"articles,omitempty"`
}

func (n *News) Path() string {
	return n.User.Dir() + "/news"
}

func (n *News) Dir() string {
	return n.Path() + "/" + n.ID.String()
}

func (n *News) Key() string {
	return n.Dir() + "/_.json"
}

func (n *News) Create(b []byte) (*News, error) {

	var v News
	if err := json.Unmarshal(b, &v); err != nil {
		log.Err(err).Msg("failed to unmarshal news")
		return nil, err
	}
	v.User = n.User
	v.ID = NewUlid()

	if err := em.Save(&v); err != nil {
		log.Err(err).Msg("failed to save news")
		return nil, err
	}

	return &v, nil
}

func (n *News) Find() error {
	return em.Find(n)
}

func (n *News) Delete() error {

	if err := em.Delete(n); err != nil {
		return err
	}

	j, err := n.job()
	if err != nil {
		return err
	}

	return em.Delete(j)
}

func (n *News) FindAll() ([]*News, error) {

	out, err := em.FindAll(n, regexp.MustCompile(`.*/news/([A-Za-z0-9]{26}/_.json)$`))
	if err != nil {
		return nil, err
	}

	for i, v := range out {
		out[i].Articles, _ = NewArticle(n.User, v.ID.String()).FindAll()
	}

	return out, nil
}

func (n *News) Work() {

	var errs error

	q := strings.Join(n.Keywords, "+")
	urls := []string{
		fmt.Sprintf("https://news.google.com/rss/search?q=%s&hl=en-US&gl=US&ceid=US:en", q),
		fmt.Sprintf("https://www.bing.com/news/search?format=rss&q=%s", q),
		fmt.Sprintf("https://www.bing.com/search?format=rss&q=%s", q),
	}

	var items []Item
	for _, u := range urls {

		rss, err := NewRSS(u)
		if err != nil {
			log.Warn().Err(err).Msg("failed to parse rss feed")
			errs = errors.Join(errs, err)
			continue
		}
		if rss.Channel.Items != nil && len(rss.Channel.Items) > 0 {
			items = append(items, rss.Channel.Items...)
		}
	}

	var article Article
	if err := article.FindLast(); err != nil {
		log.Err(err).Msg("failed to find last article, must return")
		errs = errors.Join(errs, err)
		return
	}

	for _, i := range items {
		if i.Time.UnixMilli() > article.Time {
			i.Scrub()
			if err := em.Save(&Article{
				News:   n,
				URL:    i.URL,
				Title:  i.Title,
				Source: i.Source.Value,
				Time:   i.Time.UnixMilli(),
			}); err != nil {
				errs = errors.Join(errs, err)
			}
		}
	}

	job, err := n.job()
	if err != nil {
		return
	}

	job.Results[time.Now().UTC().UnixMilli()] = errs
	em.Save(job)

	return
}

func (n *News) job() (*Job, error) {
	job, err := NewJob(n.User, n.ID).Find()
	if err != nil {
		log.Warn().Err(err).Msg("failed to find job")
		return nil, err
	}
	return job, nil
}

func NewNews(args ...any) *News {
	var v = new(News)
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
