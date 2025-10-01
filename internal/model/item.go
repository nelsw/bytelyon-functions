package model

import (
	"bytelyon-functions/internal/client/s3"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var bingRegexp = regexp.MustCompile("</?News(:\\w+)>")

type RSS struct {
	Channel struct {
		Items Items `xml:"item"`
	} `xml:"channel"`
}

type Items []Item

func MakeNewsItems(url string) (Items, error) {
	res, e := http.Get(url)
	if e != nil {
		return nil, errors.Join(errors.New("failed to http.Get url"), e)
	}
	defer res.Body.Close()

	var b []byte
	if b, e = io.ReadAll(res.Body); e != nil {
		return nil, errors.Join(errors.New("failed to io.ReadAll response body"), e)
	}

	if strings.Contains(url, "https://www.bing.com/") {
		str := bingRegexp.ReplaceAllStringFunc(string(b), func(s string) string {
			return strings.ReplaceAll(s, ":", "_")
		})
		b = []byte(str)
	}

	var rss RSS
	if e = xml.Unmarshal(b, &rss); e != nil {
		return nil, errors.Join(errors.New("failed to unmarshal xml"), e)
	}

	return rss.Channel.Items, nil
}

type Item struct {
	ID      string    `json:"id" xml:"-"`
	Created time.Time `json:"created" xml:"-"`
	URL     string    `json:"link" xml:"link"`
	Title   string    `json:"title" xml:"title"`
	Time    *DateTime `json:"published" xml:"pubDate"`
	Source  *struct {
		URL   string `json:"url,omitempty" xml:"url,attr"`
		Value string `json:"value,omitempty" xml:",chardata"`
	} `json:"source,omitempty" xml:"source"`
	NewsSource         string `json:"news_source,omitempty" xml:"News_Source"`
	NewsImage          string `json:"news_image,omitempty" xml:"News_Image"`
	NewsImageSize      string `json:"news_image_size,omitempty" xml:"News_ImageSize"`
	NewsImageMaxWidth  int    `json:"news_image_max_width,omitempty" xml:"News_ImageMaxWidth"`
	NewsImageMaxHeight int    `json:"news_image_max_height,omitempty" xml:"News_ImageMaxHeight"`
}

func (i Item) Create(db s3.Client, jobKey string) error {
	if strings.HasPrefix(i.URL, "https://news.google.com/") {
		if s, e := decodeGoogleNewsURL(i.URL); e == nil {
			i.URL = s
		}
	} else if strings.HasPrefix(i.URL, "http://www.bing.com/news/") {
		if s, e := decodeBingNewsURL(i.URL); e == nil {
			i.URL = s
		}
	}
	i.ID = base64.URLEncoding.EncodeToString([]byte(i.URL))
	i.Created = time.Now().UTC()
	return db.Save(ItemKey(jobKey, i.ID), i)
}

func ItemKey(jobKey string, itemID string) string {
	return fmt.Sprintf("%s/%s", ItemPath(jobKey), itemID)
}

func ItemPath(jobKey string) string {
	return fmt.Sprintf("%s/item", jobKey)
}

type DateTime time.Time

func (v *DateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	s = strings.Trim(s, `"`) // Remove quotes from the JSON string
	if s == "" || s == "null" {
		return nil // Handle empty or null strings
	}

	t, err := time.Parse(time.RFC1123, s) // Parse using your custom format
	if err != nil {
		return err
	}

	*v = DateTime(t)
	return nil
}

func (v *DateTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = strings.Trim(s, `"`) // Remove quotes from the JSON string
	if s == "" || s == "null" {
		return nil // Handle empty or null strings
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	*v = DateTime(t)
	return nil
}

func (v *DateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(*v).Format(time.RFC3339) + `"`), nil
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
