package model

import (
	"bytelyon-functions/pkg/db"
	"bytelyon-functions/pkg/util"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"sort"
	"strings"
	"sync"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

type rss struct {
	Channel struct {
		Items []*item `xml:"item"`
	} `xml:"channel"`
}

type item struct {
	ID     ulid.ULID `json:"id"`
	URL    string    `json:"url" xml:"link"`
	Title  string    `json:"title" xml:"title"`
	Source string    `json:"source"`
	Time   *DateTime `json:"-" xml:"pubDate"`
	Src    *struct {
		Value string `json:"value" xml:",chardata"`
	} `json:"-" xml:"source"`
}

type ProwlerNews struct {
	*Prowler
}

func NewProwlNews(p *Prowler) *ProwlerNews {
	return &ProwlerNews{p}
}

func (p *ProwlerNews) FindAll() ([]*Node, error) {

	keys, err := db.NewS3().Keys(p.Dir())
	if err != nil {
		return nil, err
	}

	var rootMap = make(map[string]*Prowler)
	var leafMap = make(map[string]map[string][]*item)

	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, key := range keys {
		k := key[:strings.LastIndex(key, "/")]
		if _, ok := leafMap[k]; !ok {
			leafMap[k] = make(map[string][]*item)
		}

		wg.Go(func() {
			b, _ := db.NewS3().Get(key)

			if strings.HasSuffix(key, "_.json") {
				var v Prowler
				_ = json.Unmarshal(b, &v)
				rootMap[k] = &v
				return
			}

			var v item
			_ = json.Unmarshal(b, &v)

			date := v.ID.Timestamp().Format("2006-01-02")

			mu.Lock()
			leafMap[k][date] = append(leafMap[k][date], &v)
			mu.Unlock()
		})
	}
	wg.Wait()
	var nodes []*Node

	for id, prowler := range rootMap {

		rootLabel := id[strings.LastIndex(id, "/")+1:]
		rootID := "news/" + rootLabel
		root := NewNode(rootID, rootLabel)
		root.Data = prowler
		nodes = append(nodes, root)

		dates := slices.Collect(maps.Keys(leafMap[id]))
		slices.Sort(dates)
		slices.Reverse(dates)

		for _, date := range dates {
			dateLabel := date
			dateID := rootID + "/" + dateLabel
			branch := NewNode(dateID, dateLabel)
			leaves := leafMap[id][date]
			branch.Data = leaves
			root.Children = append(root.Children, branch)
		}
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Label < nodes[j].Label })
	return nodes, nil
}

func (p *ProwlerNews) Go() ulid.ULID {
	return p.goSaveItems(db.NewS3(), p.goFindItems())
}

func (p *ProwlerNews) rss(s string) []*item {
	res, err := http.Get(s)
	if err != nil {
		log.Warn().Err(err).Msg("ProwlerNews - Failed to fetch rss feed")
		return nil
	}
	defer res.Body.Close()

	var b []byte
	if b, err = io.ReadAll(res.Body); err != nil {
		log.Warn().Err(err).Msg("Prowler - Failed to read rss feed")
	}

	var r rss
	if err = xml.Unmarshal(b, &r); err != nil {
		log.Warn().Err(err).Str("url", s).Msg("ProwlerNews - Failed to unmarshal rss feed")
	}
	return r.Channel.Items
}

func (p *ProwlerNews) goFindItems() []*item {
	q := strings.ReplaceAll(p.ID, ` `, `+`)
	urls := []string{
		fmt.Sprintf("https://news.google.com/rss/search?q=%s&hl=en-US&gl=US&ceid=US:en", q),
	}

	var wg sync.WaitGroup
	var items []*item
	for _, u := range urls {
		wg.Go(func() { items = append(items, p.rss(u)...) })
	}
	wg.Wait()

	sort.Slice(items, func(i, j int) bool {
		return items[i].Time.UnixMilli() < items[j].Time.UnixMilli()
	})
	log.Debug().Int("count", len(items)).Msg("ProwlerNews - Found items")
	return items
}

func (p *ProwlerNews) goSaveItems(s3 db.S3, items []*item) ulid.ULID {

	var prowled ulid.ULID
	var saved int
	var wg sync.WaitGroup
	for _, i := range items {

		wg.Go(func() {
			if p.Prowled.Timestamp().UnixMilli() > i.Time.UnixMilli() {
				return
			}

			if s, err := decodeGoogleNewsURL(i.URL); err != nil {
				log.Warn().Err(err).Msg("ProwlerNews - failed to decode google news url")
			} else {
				i.URL = s
			}

			i.ID = i.Time.ULID()

			if i.Src != nil {
				i.Source = i.Src.Value
			} else {
				i.Source = util.Domain(i.URL)
			}

			if idx := strings.LastIndex(i.Title, " - "); idx != -1 {
				i.Title = i.Title[:idx]
			}

			b, err := json.Marshal(i)
			if err != nil {
				log.Warn().Err(err).Msg("ProwlerNews - Failed to marshal article")
				return
			}

			key := fmt.Sprintf("%s%s.json", p.Prowler.Dir(), i.ID)
			if err = s3.Put(key, b); err != nil {
				log.Warn().Err(err).Msg("ProwlerNews - Failed to save article")
				return
			}
			saved++
			prowled = i.ID
		})
	}
	wg.Wait()
	log.Debug().Int("count", saved).Msg("ProwlerNews - Saved items")
	return prowled
}

func decodeGoogleNewsURL(s string) (string, error) {
	res, err := http.Get(s)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var doc *html.Node
	if doc, err = html.Parse(res.Body); err != nil {
		return "", err
	}

	// Regular expression to extract the Base64 encoded part
	encodedText := regexp.MustCompile(`/articles/(?P<encoded_url>[^?]+)`).FindStringSubmatch(s)[1]

	var fn func(*html.Node) (string, error)
	fn = func(n *html.Node) (string, error) {

		if n.Type == html.ElementNode && n.Data == "c-wiz" {

			var sg, ts string
			if e := n.FirstChild; e != nil {
				for _, att := range e.Attr {
					if att.Key == "data-n-a-sg" {
						sg = att.Val
					} else if att.Key == "data-n-a-ts" {
						ts = att.Val
					}
				}
			}
			return decodeGoogleNewsURLParts(sg, ts, encodedText)
		}

		// continue traversing every sibling per child. give em noogies.
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if u, e := fn(c); u != "" && e == nil {
				return u, nil
			}
		}
		return "", nil
	}

	return fn(doc)
}

func decodeGoogleNewsURLParts(signature, timestamp, base64Str string) (string, error) {
	endpoint := "https://news.google.com/_/DotsSplashUi/data/batchexecute"
	payload := []interface{}{
		"Fbv4je",
		fmt.Sprintf("[\"garturlreq\",[[\"X\",\"X\",[\"X\",\"X\"],null,null,1,1,\"US:en\",null,1,null,null,null,null,null,0,1],\"X\",\"X\",1,[1,1,1],1,1,null,0,0,null,0],\"%s\",%s,\"%s\"]", base64Str, timestamp, signature),
	}
	outer := [][]interface{}{payload}
	bodyBytes, _ := json.Marshal([][][]interface{}{outer})
	form := url.Values{}
	form.Set("f.req", url.QueryEscape(string(bodyBytes)))

	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString("f.req="+string(url.QueryEscape(string(bodyBytes)))))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")

	client := &http.Client{}
	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var b []byte
	if b, err = io.ReadAll(resp.Body); err != nil {
		return "", err
	}

	s := string(b)
	parts := strings.Split(s, "\n\n")
	if len(parts) < 2 {
		return "", errors.New("unexpected batchexecute response format")
	}

	payload = []interface{}{}
	if err = json.Unmarshal([]byte(parts[1]), &payload); err != nil {
		return "", err
	} else if len(payload) == 0 {
		return "", errors.New("empty payload")
	}

	entry, ok := payload[0].([]interface{})
	if !ok || len(entry) < 3 {
		return "", errors.New("unexpected entry structure")
	}

	var inner []interface{}
	if s, ok = entry[2].(string); !ok {
		return "", errors.New("missing inner json string")
	} else if err = json.Unmarshal([]byte(s), &inner); err != nil {
		return "", err
	} else if len(inner) < 2 {
		return "", errors.New("unexpected inner array")
	} else if s, ok = inner[1].(string); !ok {
		return "", errors.New("decoded url not string")
	}

	return s, nil
}
