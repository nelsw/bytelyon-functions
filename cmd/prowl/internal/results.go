package internal

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

type Results struct {
	Content string
	Doc     *goquery.Document
	Data    map[ResultType][]*Result
}

func NewResults(content string) (*Results, error) {
	var r = new(Results)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse html")
		return nil, err
	}
	r.Doc = doc
	r.Content = content
	r.Data = map[ResultType][]*Result{
		OrganicResultType:   {},
		SponsoredResultType: {},
		VideoResultType:     {},
		ForumResultType:     {},
		ArticleResultType:   {},
	}
	r.fillOrganicData()
	r.fillSponsoredData()
	return r, nil
}

func (r *Results) fillOrganicData() {
	content := r.Content
	left := strings.Index(content, `var m={`) + 7
	content = content[left:]
	content = content[:strings.Index(content, "};")]

	var vals []string
	for i, chunk := range strings.Split(content, `:[`) {

		if i == 0 {
			continue
		}

		idx := strings.LastIndex(chunk, `,`)
		if idx == -1 {
			continue
		}

		key := chunk[idx+1:]
		val := `[` + chunk[:idx]
		_, err := strconv.Atoi(key[len(key)-2 : len(key)-1])

		if err == nil && strings.Contains(val, "Source: ") {
			vals = append(vals, val)
		}
	}

	for pos := 0; pos < len(vals); pos++ {

		var o = new(Result)
		o.Position = pos

		val := strings.ReplaceAll(vals[pos], "null,", "")
		for i, v := range strings.Split(val, ",[") {

			if i == 0 {
				v = strings.ReplaceAll(v, "\"", "")
				v, _ = strconv.Unquote("\"" + v + "\"")
				o.Link = v[1:]
			} else if i == 2 {
				for i, v = range strings.Split(v, "\",\"") {
					switch i {
					case 0:
						v = strings.ReplaceAll(v, "\\u003d", "=")
						v = strings.ReplaceAll(v, "\\u0026", "&")
						v = strings.ReplaceAll(v, "\\", "")
						v = strings.ReplaceAll(v, "\"", "")
						o.Title = v
					case 1:
						v, _ = strconv.Unquote("\"" + v + "\"")
						o.Snippet = v
					case 2:
						o.Source = v
					}
				}
				break
			}
		}
		if strings.Contains(val, "WEB_RESULT_INNER") {
			o.Position = len(r.Data[OrganicResultType])
			r.Data[OrganicResultType] = append(r.Data[OrganicResultType], o)
		} else if strings.Contains(val, "COMMUNITY_MODE_WEB_RESULT") {
			o.Position = len(r.Data[ForumResultType])
			r.Data[ForumResultType] = append(r.Data[ForumResultType], o)
		} else if strings.Contains(val, "VIDEO_RESULT") {
			o.Position = len(r.Data[VideoResultType])
			r.Data[VideoResultType] = append(r.Data[VideoResultType], o)
		} else if strings.Contains(val, "NEWS_ARTICLE_RESULT") {
			o.Position = len(r.Data[ArticleResultType])
			r.Data[ArticleResultType] = append(r.Data[ArticleResultType], o)
		}
	}
}

func (r *Results) fillSponsoredData() {
	var ids []string
	r.Doc.Find(`div`).Each(func(i int, sel *goquery.Selection) {
		if _, ok := sel.Attr("data-merchant-id"); !ok {
			return
		}
		if id, ok := sel.Attr("id"); ok && id[0] == '_' {
			ids = append(ids, id)
		}
	})

	var frags []string
	for _, id := range ids {
		if left := strings.Index(r.Content, id+`',`); left > 0 {
			left += len(id) + 4
			right := strings.Index(r.Content[left:], `);})();`)
			frags = append(frags, r.Content[left:left+right])
		}
	}

	for i := range frags {
		frags[i] = strings.ReplaceAll(frags[i], "x26", "&")
		frags[i] = strings.ReplaceAll(frags[i], "x27", "'")
		frags[i] = strings.ReplaceAll(frags[i], "xb2", "²")
		frags[i] = strings.ReplaceAll(frags[i], "x3d", "=")
		frags[i] = strings.ReplaceAll(frags[i], "x22", "")
		frags[i] = strings.ReplaceAll(frags[i], "x3c", "<")
		frags[i] = strings.ReplaceAll(frags[i], "x3e", ">")
		frags[i] = strings.ReplaceAll(frags[i], "&amp;", "&")
		frags[i] = strings.ReplaceAll(frags[i], `\`, ``)
	}

	for pos, f := range frags {
		var result = new(Result)
		result.Position = pos

		d, err := html.Parse(strings.NewReader(f))
		if err != nil {
			log.Warn().Err(err).Msg("failed to parse sponsored html")
			continue
		}

		gd := goquery.NewDocumentFromNode(d)

		goquery.NewDocumentFromNode(d).Find(`span`).Each(func(i int, sel *goquery.Selection) {

			t := strings.TrimSpace(sel.Text())
			if len(t) == 0 || t[0] != '$' {
				return
			}

			t = strings.ReplaceAll(t, ",", "")[1:]
			if price, e := strconv.ParseFloat(t, 64); e == nil {
				result.Price = price
			}
		})

		gd.Find(`div`).Each(func(i int, sel *goquery.Selection) {

			t := strings.TrimSpace(sel.Text())

			if _, ok := sel.Attr("aria-label"); ok {
				result.Source = t
				return
			}

			if val, ok := sel.Attr("role"); ok && val == "heading" {
				result.Title = t
				return
			}

		})

		gd.Find(`a`).Each(func(i int, sel *goquery.Selection) {
			if result.Link != "" {
				return
			}
			if val, ok := sel.Attr("href"); ok && strings.Contains(val, "https://") {
				result.Link = val
				return
			}
		})

		r.Data[SponsoredResultType] = append(r.Data[SponsoredResultType], result)
	}
}
