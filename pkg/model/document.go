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

	var attempts int
	if len(i) > 0 {
		attempts = i[0]
	}

	if attempts > 6 {
		return errors.New("max fetch attempts reached")
	} else if attempts > 0 {
		time.Sleep(time.Second * time.Duration(i[0]*10))
	}

	res, err := http.Get(d.URL)
	if err != nil {
		log.Err(err).Str("URL", d.URL).Msg("Document get")
		return d.Fetch(attempts + 1)
	}

	defer res.Body.Close()

	if d.Node, err = html.Parse(res.Body); err != nil {
		log.Err(err).Str("URL", d.URL).Msg("Document parse")
		return err
	}

	return nil
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
