package model

import (
	"bytelyon-functions/pkg/db"
	"bytelyon-functions/pkg/util"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
)

var (
	googleSearchInputSelectors = []string{
		"input[name='q']",
		"input[title='Search']",
		"input[aria-label='Search']",
		"textarea[title='Search']",
		"textarea[name='q']",
		"textarea[aria-label='Search']",
		"textarea",
	}
)

type ProwlerSearch struct {
	ID       ulid.ULID `json:"id"`
	db.S3    `json:"-"`
	*PW      `json:"-"`
	*Prowler `json:"-"`
}

func NewProwlSearch(p *Prowler) *ProwlerSearch {
	return &ProwlerSearch{
		ID:      NewUlid(),
		S3:      db.NewS3(),
		Prowler: p,
	}
}

func (p *ProwlerSearch) Dir() string {
	return p.Prowler.Dir() + p.ID.String() + "/"
}

func (p *ProwlerSearch) FindAll() ([]*Node, error) {
	keys, err := db.NewS3().Keys(p.Prowler.Dir())
	if err != nil {
		return nil, err
	}

	var jsonKeys []string
	for _, key := range keys {
		if strings.HasSuffix(key, ".json") {
			jsonKeys = append(jsonKeys, key)
		}
	}

	var fullKeyPrefix = "user/" + p.UserID.String() + "/prowler/"

	var nodeIds []string
	for _, k := range jsonKeys {
		id := strings.TrimPrefix(k, fullKeyPrefix)
		nodeIds = append(nodeIds, id)
	}

	var rootMap = make(map[string]*Node)
	for _, nodeId := range nodeIds {
		if strings.HasSuffix(nodeId, "/_.json") {
			b, _ := p.S3.Get(fullKeyPrefix + nodeId)
			var pp Prowler
			_ = json.Unmarshal(b, &pp)

			id := strings.TrimSuffix(nodeId, "/_.json")
			label := strings.TrimPrefix(id, "search/")
			node := NewNode(id, label, &pp)
			rootMap[id] = node
		}
	}

	var serpKeys []string
	for _, k := range jsonKeys {
		if strings.Contains(k, "/serp/") {
			serpKey := strings.TrimPrefix(k, fullKeyPrefix)

			serpKeys = append(serpKeys, serpKey)
		}
	}

	var dateMap = make(map[string][]*Node)
	for _, serpKey := range serpKeys {

		str := strings.TrimPrefix(serpKey, "search/")
		idx := strings.Index(str, "/")
		rootNodeId := "search/" + str[:idx]

		str = str[idx+1:]
		idx = strings.Index(str, "/")
		ulidStr := str[:idx]

		dateTime := ulid.MustParse(ulidStr).Timestamp().Format(time.DateTime)

		b, _ := p.S3.Get(fullKeyPrefix + serpKey)
		var data map[string]any
		_ = json.Unmarshal(b, &data)

		img, _ := p.S3.URL(strings.ReplaceAll(fullKeyPrefix+serpKey, ".json", ".png"), 30)
		html, _ := p.S3.URL(strings.ReplaceAll(fullKeyPrefix+serpKey, ".json", ".html"), 30)

		idx = strings.LastIndex(serpKey, "/")
		id := serpKey[idx+1:]
		id = strings.TrimSuffix(id, ".json")
		node := NewNode(rootNodeId+"/"+dateTime, dateTime, map[string]any{
			"id":   id,
			"img":  img,
			"html": html,
			"json": data,
		})

		dateMap[rootNodeId] = append(dateMap[rootNodeId], node)
	}

	var rootIds = slices.Collect(maps.Keys(rootMap))
	sort.Strings(rootIds)

	for _, rootId := range rootIds {
		children := dateMap[rootId]
		sort.Slice(children, func(i, j int) bool {
			iTime, _ := time.Parse(time.DateTime, children[i].Label)
			jTime, _ := time.Parse(time.DateTime, children[j].Label)
			return iTime.UnixMilli() < jTime.UnixMilli()
		})
		rootMap[rootId].Children = children
	}

	var nodes []*Node

	for _, rootId := range rootIds {
		nodes = append(nodes, rootMap[rootId])
	}

	return nodes, nil
}

func (p *ProwlerSearch) Go() ulid.ULID {
	var fn func(bool) ulid.ULID

	fn = func(headless bool) ulid.ULID {

		var err error
		p.PW, err = NewPW(headless)
		if err != nil {
			log.Warn().Err(err).Msg("ProwlerSearch - Failed to initialize PW")
			return ulid.Zero
		}
		defer p.PW.Close()

		log.Info().Bool("headless", headless).Msg("ProwlerSearch - Working ... ")

		var prowled ulid.ULID
		if prowled, err = p.worker(); err != nil && headless {
			log.Warn().Err(err).Msg("ProwlerSearch - Headless Failed; retrying with head ...")
			return fn(false)
		}

		if err != nil {
			log.Warn().Err(err).Bool("headless", headless).Msg("ProwlerSearch - Failed!")
			return ulid.Zero
		}

		log.Info().Bool("headless", headless).Msg("ProwlerSearch - Success!")
		return prowled
	}

	return fn(true)
}

func (p *ProwlerSearch) worker() (prowled ulid.ULID, err error) {

	defer p.PW.Close()

	var page playwright.Page
	if page, err = p.PW.NewPage(); err != nil {
		return
	}
	defer page.Close()

	var res playwright.Response
	if res, err = p.PW.GoTo(page, "https://www.google.com"); err != nil {
		return
	} else if err = p.PW.IsBlocked(page, res); err != nil {
		return
	} else if err = p.PW.Click(page, googleSearchInputSelectors...); err != nil {
		return
	} else if err = p.PW.Type(page, p.Prowler.ID); err != nil {
		return
	} else if err = p.PW.Press(page, "Enter"); err != nil {
		return
	} else if err = p.PW.WaitForLoadState(page); err != nil {
		return
	} else if err = p.PW.IsBlocked(page); err != nil {
		return
	}

	prowled = p.save(page)

	targetCount := len(p.Targets)
	log.Info().Msgf("ProwlerSearch - Pages [%d]", targetCount)

	if targetCount == 0 {
		return
	}

	var locators []playwright.Locator
	if locators, err = page.Locator(fmt.Sprintf(`[data-dtld]`), playwright.PageLocatorOptions{}).All(); err != nil {
		return
	}

	if len(locators) == 0 {
		log.Warn().Msg("ProwlerSearch - No Target Locators Found")
		return
	}

	var att string
	for _, l := range locators {

		if att, err = l.GetAttribute("data-dtld", playwright.LocatorGetAttributeOptions{
			Timeout: util.Ptr(5_000.0),
		}); err != nil {
			log.Warn().Err(err).Msg("ProwlerSearch - Failed to get Target Locator Attribute")
			continue
		}

		log.Debug().Str("found", att).Msg("ProwlerSearch - Locator")
		if !p.Targets.Follow(att) {
			continue
		}

		log.Info().Msgf("ProwlerSearch - Target Found [%s]", att)

		if err = l.Click(playwright.LocatorClickOptions{Timeout: util.Ptr(5_000.0)}); err != nil {
			log.Warn().Err(err).Msg("ProwlerSearch - Failed to Click Target Locator")
			continue
		}

		targetPage, pageErr := p.PW.NewPage(func() error { return l.Click() })
		if pageErr != nil {
			log.Warn().Err(pageErr).Msg("ProwlerSearch - Failed to Click Target")
		} else {
			prowled = p.save(targetPage)
			targetPage.Close()
		}
	}

	return prowled, nil
}

func (p *ProwlerSearch) save(page playwright.Page) (prowled ulid.ULID) {

	prowled = NewUlid()
	url := page.URL()
	domain := util.Domain(url)

	path := p.Dir()
	if domain == "google.com" {
		path += "serp/"
	} else {
		path += "target/"
	}
	path += prowled.String()

	var err error

	var img []byte
	if img, err = page.Screenshot(playwright.PageScreenshotOptions{FullPage: util.Ptr(true)}); err != nil {
		log.Warn().Err(err).Str("domain", domain).Msg("PW - Failed to Screenshot Page")
	} else {
		p.S3.Put(path+".png", img)
	}

	var content string
	if content, err = page.Content(); err != nil {
		log.Warn().Err(err).Str("domain", domain).Msg("PW - Failed to get Page Content")
	} else {
		p.S3.Put(path+".html", []byte(content))
	}

	var title string
	if title, err = page.Title(); err != nil {
		log.Warn().Err(err).Str("domain", domain).Msg("PW - Failed to get Page Title")
	}

	m := map[string]any{
		"url":    url,
		"domain": domain,
		"title":  title,
	}

	if m["domain"] == "google.com" {
		m["results"] = p.PW.Data(url, content)
		m["targets"] = p.Targets
	}

	var data []byte
	if data, err = json.Marshal(&m); err != nil {
		log.Warn().Err(err).Str("domain", domain).Msg("PW - Failed to Marshal Page")
	} else {
		p.S3.Put(path+".json", data)
	}

	log.Info().
		Str("domain", domain).
		Msg("PW - Saved Page")

	return
}
