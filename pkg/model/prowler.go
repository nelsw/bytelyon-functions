package model

import (
	"bytelyon-functions/pkg/db"
	"bytelyon-functions/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
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
	Prowled   ulid.ULID     `json:"prowled,omitempty"`
	Frequency time.Duration `json:"frequency"`
	Duration  time.Duration `json:"duration"`
	Targets   Targets       `json:"targets,omitempty"`
	Sessions  []any         `json:"sessions,omitempty"`
	Disabled  bool          `json:"disabled"`
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

func (p *Prowler) FindAll() (any, error) {

	switch p.Type {
	case SearchProwlerType:
		return NewProwlSearch(p).FindAll()
	case SitemapProwlerType:
		return NewProwlSitemap(p).FindAll()
	case NewsProwlerType:
		pp := &ProwlerNews{p}
		return pp.FindAll()
	}

	return nil, errors.New("not implemented")
}

func (p *Prowler) findAllSearches() ([]*Prowler, error) {

	s3 := db.NewS3()

	keys, err := s3.Keys(p.Dir())
	if err != nil {
		return nil, err
	} else if len(keys) == 0 {
		return nil, errors.New("no keys found")
	}

	var uniqueJsonNames = make(map[string][]string)
	for _, k := range keys {
		if !strings.HasSuffix(k, ".json") {
			continue
		}
		k = strings.TrimPrefix(k, p.Dir())
		idx := strings.Index(k, "/")
		key := p.Dir() + k[:idx]
		val := k[idx+1:]
		uniqueJsonNames[key] = append(uniqueJsonNames[key], val)
	}

	var sessionGroups = make(map[string][]string)
	for k, v := range uniqueJsonNames {
		for _, vv := range v {
			idx := strings.Index(vv, "/")
			if idx == -1 {
				continue
			}
			sessionId := k + "/" + vv[:idx]
			sessionId = strings.TrimPrefix(sessionId, p.Dir())
			sessionGroups[sessionId] = append(sessionGroups[sessionId], vv[idx:])
		}
	}

	var searchGroups = make(map[string][]string)
	for k, v := range sessionGroups {
		idx := strings.Index(k, "/")
		key := k[:idx]
		sessionId := k[idx+1:]
		for _, vv := range v {
			searchGroups[key] = append(searchGroups[key], sessionId+vv)
		}
	}

	type Page struct {
		Data  any     `json:"data"`
		Img   string  `json:"img"`
		Html  string  `json:"html"`
		Pages []*Page `json:"pages,omitempty"`
	}

	var prowlers []*Prowler
	for k, v := range searchGroups {

		k = p.Dir() + k
		fmt.Println("key", k)

		var x = new(Prowler)
		b, _ := s3.Get(k + "/_.json")
		_ = json.Unmarshal(b, x)

		for _, vv := range v {

			key := k + "/" + vv
			b, _ = s3.Get(key)
			var data map[string]any
			_ = json.Unmarshal(b, &data)

			img, _ := s3.URL(strings.ReplaceAll(key, ".json", ".png"), 30)
			html, _ := s3.URL(strings.ReplaceAll(key, ".json", ".html"), 30)

			page := &Page{Img: img, Html: html, Data: &data}
			if strings.Contains(vv, "serp") {
				x.Sessions = append(x.Sessions, page)
			} else {
				x.Sessions[len(x.Sessions)-1].(*Page).Pages = append(x.Sessions[len(x.Sessions)-1].(*Page).Pages, page)
			}
		}
		prowlers = append(prowlers, x)
	}

	//var prowlersMap = make(map[string]*Prowler)
	//idPageFileMap := make(map[string]map[string]map[string]int)

	//for _, k := range keys {
	//	fmt.Println(k)
	//	if strings.HasSuffix(k, "/_.json") {
	//
	//		b, _ := s3.Get(k)
	//		var e *Prowler
	//		_ = json.Unmarshal(b, &e)
	//		prowlersMap[e.ID] = e
	//		prowlers = append(prowlers, e)
	//
	//	} else if loadResults {
	//
	//		prowler := prowlers[len(prowlers)-1]
	//
	//		if _, ok := idPageFileMap[prowler.ID]; !ok {
	//			idPageFileMap[prowler.ID] = make(map[string]map[string]int)
	//		}
	//
	//		ids := strings.Split(strings.TrimPrefix(k, prowler.Dir()), "/")
	//		page := ids[0]
	//		file := ids[1]
	//		if idx := strings.Index(file, "."); idx > 0 {
	//			file = file[:idx]
	//		}
	//
	//		if _, ok := idPageFileMap[prowler.ID][page]; !ok {
	//			idPageFileMap[prowler.ID][page] = make(map[string]int)
	//		}
	//
	//		idPageFileMap[prowler.ID][page][file]++
	//	}
	//}

	//type Result struct {
	//	Data     any      `json:"data"`
	//	Img      string   `json:"img"`
	//	Html     string   `json:"html"`
	//	Followed []Result `json:"followed,omitempty"`
	//}
	//
	//var r Result
	//for id, m1 := range idPageFileMap {
	//	for page, m2 := range m1 {
	//		for file, count := range m2 {
	//			img, _ := s3.URL(p.Dir()+fmt.Sprintf("%s/%s/%s.png", id, page, file), 30)
	//			html, _ := s3.URL(p.Dir()+fmt.Sprintf("%s/%s/%s.html", id, page, file), 30)
	//			if count < 3 {
	//				r.Followed = append(r.Followed, Result{Img: img, Html: html})
	//				continue
	//			}
	//			r.Img = img
	//			r.Html = html
	//			b, _ := s3.Get(fmt.Sprintf("%s%s/%s/%s.json", p.Dir(), id, page, file))
	//			var m map[string]any
	//			_ = json.Unmarshal(b, &m)
	//			r.Data = m
	//		}
	//	}
	//	prowlersMap[id].Sessions = append(prowlersMap[id].Sessions, r)
	//}

	return prowlers, nil
}

func (p *Prowler) findAllSimpleTypes() ([]*Prowler, error) {

	s3 := db.NewS3()

	keys, err := s3.Keys(p.Dir())
	if err != nil {
		return nil, err
	} else if len(keys) == 0 {
		return nil, errors.New("no keys found")
	}

	var prowlers []*Prowler

	var m = make(map[int][]any)
	for _, k := range keys {
		b, _ := s3.Get(k)
		if !strings.HasSuffix(k, "_.json") {
			var a any
			_ = json.Unmarshal(b, &a)
			m[len(prowlers)] = append(m[len(prowlers)], a)
			continue
		}
		var prowler = new(Prowler)
		_ = json.Unmarshal(b, prowler)
		prowler.Sessions = append(prowler.Sessions, m[len(prowlers)]...)
		prowlers = append(prowlers, prowler)
	}

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
	var prowled ulid.ULID

	switch p.Type {
	case SearchProwlerType:
		prowled = NewProwlSearch(p).Go()
	case SitemapProwlerType:
		prowled = NewProwlSitemap(p).Go()
	case NewsProwlerType:
		prowled = NewProwlNews(p).Go()
	}

	if prowled.IsZero() {
		log.Warn().Msgf("Prowl - Empty Prowl [%s/%s]", p.Type, p.ID)
		return
	}

	p.Duration = time.Since(ts)
	if err := db.Save(p); err != nil {
		log.Warn().Err(err).Msgf("Prowl - Failed to save Prowler [%s/%s]", p.Type, p.ID)
		return
	}

	log.Info().Msgf("Prowl - Prowled [%s/%s]", p.Type, p.ID)
}
