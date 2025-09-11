package model

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/s3"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"
)

type URL struct {
	URL       string    `json:"url" fake:"{url}"`
	CreatedAt time.Time `json:"created_at" fake:"{datetime}"`
}

func MakeURL(s string) URL {
	return URL{URL: s, CreatedAt: time.Now().UTC()}
}

func (u URL) Path() string {
	return "page"
}

func (u URL) Key() string {
	return fmt.Sprintf("%s/%s", u.Path(), base64.URLEncoding.EncodeToString([]byte(u.URL)))
}

func (u URL) Validate() (err error) {
	_, err = url.Parse(u.URL)
	return
}

func (u URL) Create(db s3.Client) error {
	return db.Put(u.Key(), app.MustMarshal(u))
}
