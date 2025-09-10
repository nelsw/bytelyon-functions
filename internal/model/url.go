package model

import (
	"encoding/base64"
	"fmt"
	"net/url"
)

type URL struct {
	URL string `json:"url" fake:"{url}"`
}

func (u URL) Path() string {
	return "page"
}

func (u URL) Key() string {
	return fmt.Sprintf("%s/%s/_.json", u.Path(), base64.URLEncoding.EncodeToString([]byte(u.URL)))
}

func (u URL) Validate() (err error) {
	_, err = url.Parse(u.URL)
	return
}
