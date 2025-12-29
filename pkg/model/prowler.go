package model

import (
	"bytelyon-functions/pkg/db"
	"bytelyon-functions/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type Prowler struct {
	UserID ulid.ULID `json:"user_id"`
	// ID is either a URL or Query string
	ID        string        `json:"id"`
	Type      ProwlerType   `json:"type"`
	Prowled   ulid.ULID     `json:"prowled"`
	Frequency time.Duration `json:"frequency"`
	Duration  time.Duration `json:"duration"`
	Targets   Targets       `json:"targets,omitempty"`
	Results   []any         `json:"results,omitempty"`
}

func (p *Prowler) Dir() string {
	id := p.ID
	if p.Type == SitemapProwlerType {
		id = util.Domain(p.ID)
	}
	if id == "" {
		return fmt.Sprintf("user/%s/prowler/%s/", p.UserID, p.Type)
	}
	return fmt.Sprintf("user/%s/prowler/%s/%s/", p.UserID, p.Type, id)
}

func (p *Prowler) Key() string {
	return p.Dir() + "_.json"
}

func (p *Prowler) FindAll(eager ...bool) (pp []*Prowler, err error) {

	switch p.Type {
	case SearchProwlerType:
		return p.findAllSearches(true)
	case SitemapProwlerType, NewsProwlerType:
		return p.findAllSimpleTypes(true)
	}

	return nil, errors.New("not implemented")
}

func (p *Prowler) findAllSearches(eager ...bool) ([]*Prowler, error) {

	loadResults := len(eager) > 0 && eager[0]

	s3 := db.NewS3()

	keys, err := s3.Keys(p.Dir())
	if err != nil {
		return nil, err
	} else if len(keys) == 0 {
		return nil, errors.New("no keys found")
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	var prowlers []*Prowler
	var prowlersMap = make(map[string]*Prowler)
	idPageFileMap := make(map[string]map[string]map[string]int)

	for _, k := range keys {

		if strings.HasSuffix(k, "/_.json") {

			b, _ := s3.Get(k)
			var e *Prowler
			_ = json.Unmarshal(b, &e)
			prowlersMap[e.ID] = e
			prowlers = append(prowlers, e)

		} else if loadResults {

			prowler := prowlers[len(prowlers)-1]

			if _, ok := idPageFileMap[prowler.ID]; !ok {
				idPageFileMap[prowler.ID] = make(map[string]map[string]int)
			}

			ids := strings.Split(strings.TrimPrefix(k, prowler.Dir()), "/")
			page := ids[0]
			file := ids[1]
			if idx := strings.Index(file, "."); idx > 0 {
				file = file[:idx]
			}

			if _, ok := idPageFileMap[prowler.ID][page]; !ok {
				idPageFileMap[prowler.ID][page] = make(map[string]int)
			}

			idPageFileMap[prowler.ID][page][file]++
		}
	}

	type Result struct {
		Data []byte `json:"data"`
		Img  string `json:"img"`
		Html string `json:"html"`
	}

	for id, m1 := range idPageFileMap {
		for page, m2 := range m1 {
			for file, count := range m2 {
				var r Result
				r.Img, _ = s3.URL(p.Dir()+fmt.Sprintf("%s/%s/%s.png", id, page, file), 30)
				r.Html, _ = s3.URL(p.Dir()+fmt.Sprintf("%s/%s/%s.html", id, page, file), 30)
				if count > 2 {
					r.Data, _ = s3.Get(fmt.Sprintf("%s%s/%s/%s.json", p.Dir(), id, page, file))
				}
				prowlersMap[id].Results = append(prowlersMap[id].Results, r)
			}
		}
	}

	return prowlers, nil
}

func (p *Prowler) findAllSimpleTypes(eager ...bool) ([]*Prowler, error) {

	loadResults := len(eager) > 0 && eager[0]

	s3 := db.NewS3()

	keys, err := s3.Keys(p.Dir())
	if err != nil {
		return nil, err
	} else if len(keys) == 0 {
		return nil, errors.New("no keys found")
	}

	var prowlersMap = make(map[string]*Prowler)

	var e Prowler
	var a any

	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	for _, k := range keys {
		if strings.HasSuffix(k, "/_.json") {
			b, _ := s3.Get(k)
			id := strings.TrimSuffix(k, "/_.json")
			_ = json.Unmarshal(b, &e)
			prowlersMap[id] = &e
		} else if loadResults {
			id := strings.TrimSuffix(k, ".json")
			id = id[:len(id)-27]
			b, _ := s3.Get(k)
			_ = json.Unmarshal(b, &a)
			prowlersMap[id].Results = append(prowlersMap[id].Results, a)
		}
	}
	prowlers := slices.Collect(maps.Values(prowlersMap))
	sort.Slice(prowlers, func(i, j int) bool {
		return prowlers[i].Prowled.Timestamp().UnixMilli() > prowlers[j].Prowled.Timestamp().UnixMilli()
	})

	return prowlers, nil
}

func (p *Prowler) Prowl() {

	if p.Frequency == 0 && !p.Prowled.IsZero() {
		log.Info().Msg("Prowler - Already prowled ...")
		return
	}

	if p.Prowled.Timestamp().Add(p.Frequency).After(time.Now()) {
		log.Info().Msg("Prowler - Too soon to prowl ...")
		return
	}

	ts := time.Now()

	switch p.Type {
	case SearchProwlerType:
		p.Prowled = NewProwlSearch(p).Go()
	case SitemapProwlerType:
		p.Prowled = NewProwlSitemap(p).Go()
	case NewsProwlerType:
		p.Prowled = NewProwlNews(p).Go()
	}

	p.Duration = time.Since(ts)
	if err := db.Save(p); err != nil {
		log.Warn().Err(err).Msgf("Prowl - Failed to save Prowler [%s/%s]", p.Type, p.ID)
		return
	}
	log.Info().Msgf("Prowl - Prowled [%s/%s]", p.Type, p.ID)
}
