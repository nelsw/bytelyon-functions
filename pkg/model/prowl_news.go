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
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

var (
	bingRssRegexp = regexp.MustCompile("</?News(:\\w+)>")
)

type rss struct {
	Channel struct {
		Items []*item `xml:"item"`
	} `xml:"channel"`
}

type item struct {
	ProwlerNews *ProwlNews `json:"-"`
	ID          ulid.ULID  `json:"id"`
	URL         string     `json:"url" xml:"link"`
	Title       string     `json:"title" xml:"title"`
	Source      string     `json:"source"`
	Time        *DateTime  `json:"-" xml:"pubDate"`
	Src         *struct {
		Value string `json:"value" xml:",chardata"`
	} `json:"-" xml:"source"`
}

func (i item) String() string {
	return i.ProwlerNews.String() + util.Path("item", i.ID)
}

type ProwlNews struct {
	*Prowl
}

func NewProwlNews(p *Prowl) *ProwlNews {
	return &ProwlNews{p}
}

func (p *ProwlNews) String() string {
	return p.Prowl.String()
}

func (p *ProwlNews) Go() {
	q := strings.ReplaceAll(p.Prowl.Prowler.ID, ` `, `+`)
	urls := []string{
		fmt.Sprintf("https://news.google.com/rss/search?q=%s&hl=en-US&gl=US&ceid=US:en", q),
		fmt.Sprintf("https://www.bing.com/news/search?format=rss&q=%s", q),
		fmt.Sprintf("https://www.bing.com/search?format=rss&q=%s", q),
	}

	var items []*item
	for _, u := range urls {
		items = append(items, p.rss(u)...)
	}

	for _, i := range items {

		if p.Prowl.Prowler.Prowled.Timestamp().UnixMilli() > i.Time.UnixMilli() {
			continue
		}

		if strings.HasPrefix(i.URL, "https://news.google.com/") {
			if s, err := decodeGoogleNewsURL(i.URL); err != nil {
				log.Warn().Err(err).Msg("failed to decode google news url")
			} else {
				i.URL = s
			}
		} else if strings.HasPrefix(i.URL, "http://www.bing.com/") {
			if s, err := decodeBingNewsURL(i.URL); err != nil {
				log.Warn().Err(err).Msg("failed to decode bing news url")
			} else {
				i.URL = s
			}
		}

		i.ProwlerNews = p
		i.ID = i.Time.ULID()
		if i.Src != nil {
			i.Source = i.Src.Value
		} else {
			i.Source = util.Domain(i.URL)
		}

		if err := db.Save(i); err != nil {
			log.Warn().Err(err).Msg("Prowler - Failed to save article")
		}
	}
}

func (p *ProwlNews) rss(s string) []*item {
	res, err := http.Get(s)
	if err != nil {
		log.Warn().Err(err).Msg("Prowler - Failed to fetch rss feed")
		return nil
	}
	defer res.Body.Close()

	var b []byte
	if b, err = io.ReadAll(res.Body); err != nil {
		log.Warn().Err(err).Msg("Prowler - Failed to read rss feed")
	}

	if strings.Contains(s, "https://www.bing.com/") {
		str := bingRssRegexp.ReplaceAllStringFunc(string(b), func(s string) string {
			return strings.ReplaceAll(s, ":", "_")
		})
		b = []byte(str)
	}
	var r rss
	if err = xml.Unmarshal(b, &r); err != nil {
		log.Warn().Err(err).Msg("Prowler - Failed to unmarshal rss feed")
	}
	return r.Channel.Items
}

func decodeBingNewsURL(s string) (string, error) {
	parts := strings.Split(s, "&url=")
	if len(parts) == 1 {
		return "", errors.New("invalid bing news url")
	} else if parts = strings.Split(parts[1], "&"); len(parts) == 1 {
		return "", errors.New("invalid bing news url")
	}
	return url.QueryUnescape(parts[0])
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
