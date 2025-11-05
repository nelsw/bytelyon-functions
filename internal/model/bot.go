package model

import (
	"bytelyon-functions/internal/service/playwrighter"
	"bytelyon-functions/internal/service/s3"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/playwright-community/playwright-go"
)

type BotType int

const (
	AdWordsType BotType = iota + 1
	SitemapType
)

type Bot struct {
	User   *User     `json:"-"`
	ID     ulid.ULID `json:"id"`
	Type   BotType   `json:"type"`
	Query  string    `json:"query"`
	Domain string    `json:"domain"`
	Auth   Auth      `json:"credentials"`
	URL    URL       `json:"url"`
}

func NewBot(user *User) *Bot {
	return &Bot{User: user}
}

func (b *Bot) Validate() error {

	if b.Type == AdWordsType {
		if b.Query == "" {
			return errors.New("invalid query")
		}
		if b.Domain == "" {
			return errors.New("invalid domain")
		}
		return nil
	}

	if b.Type == SitemapType {
		if err := b.URL.Validate(); err != nil {
			return err
		}
		return nil
	}

	return errors.New("invalid bot type; must be one of: 1 (AdWords), 2 (Sitemap)")
}

func (b *Bot) Path() string {
	return b.User.Path() + "/sitemap"
}

func (b *Bot) Key() string {
	return b.Path() + "/" + b.ID.String() + "/_.json"
}

func (b *Bot) FindAll(db s3.Service) ([]Bot, error) {

	keys, err := db.Keys(b.Path(), "", 1000)
	if err != nil {
		return nil, err
	}

	var vv []Bot
	for _, k := range keys {

		o, e := db.Get(k)
		if e != nil {
			err = errors.Join(err, e)
			continue
		}

		var v Bot
		if e = json.Unmarshal(o, &v); e != nil {
			err = errors.Join(err, e)
			continue
		}

		vv = append(vv, v)
	}

	return vv, err
}

func (b *Bot) Create(db s3.Service, body []byte) (*Bot, error) {

	if err := json.Unmarshal(body, b); err != nil {
		return nil, err
	} else if err = b.Validate(); err != nil {
		return nil, err
	}

	svc, err := playwrighter.New()
	if err != nil {
		return nil, err
	}

	var page playwright.Page
	page, err = svc.Goto("https://www.google.com")
	if err != nil {
		return nil, err
	} else {
		fmt.Println(page)
	}

	if body, err = json.Marshal(b); err != nil {
		return nil, err
	} else if err = db.Put(b.Key(), body); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Bot) Delete(db s3.Service, id string) (any, error) {
	b.ID, _ = ulid.Parse(id)
	return nil, db.Delete(b.Key())
}
