package model

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var bingRegexp = regexp.MustCompile("</?News(:\\w+)>")

type RSS struct {
	Channel struct {
		Items []*Item `xml:"item"`
	} `xml:"channel"`
}

func NewRSS(s string) (*RSS, error) {
	res, e := http.Get(s)
	if e != nil {
		return nil, errors.Join(errors.New("failed to http.Get url"), e)
	}
	defer res.Body.Close()

	var b []byte
	if b, e = io.ReadAll(res.Body); e != nil {
		return nil, errors.Join(errors.New("failed to io.ReadAll response body"), e)
	}

	if strings.Contains(s, "https://www.bing.com/") {
		str := bingRegexp.ReplaceAllStringFunc(string(b), func(s string) string {
			return strings.ReplaceAll(s, ":", "_")
		})
		b = []byte(str)
	}

	var rss RSS
	if e = xml.Unmarshal(b, &rss); e != nil {
		return nil, errors.Join(errors.New("failed to unmarshal xml"), e)
	}

	return &rss, nil
}
