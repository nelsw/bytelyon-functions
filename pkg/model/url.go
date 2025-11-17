package model

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

var (
	UrlExtensions       = []string{".com", ".org", ".net", ".co", ".us", ".info", ".biz", ".edu", ".gov", ".io"}
	InvalidProtocolErr  = errors.New("invalid url; must start with http:// or https://")
	InvalidExtensionErr = errors.New("invalid url; requires a valid extension: " + strings.Join(UrlExtensions, ","))
)

type URL string

func MakeURL(s string) URL {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "/")
	return URL(s)
}

func DecodeURL(s string) (URL, error) {

	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	return MakeURL(string(b)), nil
}

func (u URL) Encode() string {
	return base64.URLEncoding.EncodeToString(u.Bytes())
}

func (u URL) Validate() error {
	if !u.StartsWith("http://", "https://") {
		return InvalidProtocolErr
	}
	if !MakeURL(u.Domain()).EndsWith(UrlExtensions...) {
		return InvalidExtensionErr
	}
	return nil
}

func (u URL) Domain() string {
	s := u.String()
	s = strings.TrimPrefix(s, "http://")
	s = strings.TrimPrefix(s, "https://")
	s = strings.TrimPrefix(s, "www.")
	if i := strings.Index(s, "/"); i > 0 {
		s = s[:i]
	}
	return s
}

func (u URL) String() string {
	return string(u)
}

func (u URL) Bytes() []byte {
	return []byte(u)
}

func (u URL) Prepend(p URL) URL {
	return MakeURL(p.String() + u.String())
}

func (u URL) Append(a URL) URL {
	return MakeURL(u.String() + a.String())
}

func (u URL) StartsWith(ss ...string) bool {
	return u.has(true, ss)
}

func (u URL) EndsWith(ss ...string) bool {
	return u.has(false, ss)
}

func (u URL) has(prefix bool, ss []string) bool {
	if ss == nil || len(ss) == 0 {
		return false
	}
	for _, s := range ss {
		if prefix && strings.HasPrefix(u.String(), s) {
			return true
		}
		if !prefix && strings.HasSuffix(u.String(), s) {
			return true
		}
	}
	return false
}

func (u URL) Document() (*html.Node, error) {
	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return html.Parse(res.Body)
}
