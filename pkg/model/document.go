package model

import (
	"errors"
	"maps"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

type Document struct {
	*html.Node
	URL string
}

func NewDocument(url string) *Document {
	return &Document{URL: url}
}

func (d *Document) Fetch(i ...int) error {

	var attempt int
	if len(i) > 0 {
		attempt = i[0]
	} else {
		attempt = 1
	}

	if attempt >= 6 {
		return errors.New("max fetch attempts reached")
	}

	if attempt > 1 {
		time.Sleep(time.Second * time.Duration(attempt))
	}

	c := &http.Client{
		Timeout: time.Second * time.Duration((attempt)*3),
	}
	r, err := http.NewRequest(http.MethodGet, d.URL, nil)
	log.Err(err).Str("URL", d.URL).Msg("client request")
	if err != nil {
		return d.Fetch(attempt + 1)
	}

	r.Close = true

	var res *http.Response
	res, err = c.Do(r)
	log.Err(err).Str("URL", d.URL).Msg("client do")
	if err != nil {
		return d.Fetch(attempt + 1)
	}

	d.Node, err = html.Parse(res.Body)
	log.Err(err).Str("URL", d.URL).Msg("Document parse")

	cErr := res.Body.Close()
	log.Err(cErr).Str("URL", d.URL).Msg("Res body close")

	return err
}

func (d *Document) anchors() []string {
	m := make(map[string]bool)
	goquery.NewDocumentFromNode(d.Node).Find(`a`).Each(func(i int, sel *goquery.Selection) {
		if a, ok := sel.Attr("href"); ok {
			m[strings.TrimSpace(a)] = true
		}
	})
	return slices.Collect(maps.Keys(m))
}
