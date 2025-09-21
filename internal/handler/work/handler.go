package work

import (
	"bytelyon-functions/internal/client/s3"
	"bytelyon-functions/internal/handler/sitemap"
	"bytelyon-functions/internal/model"
	"context"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

var bingRegexp = regexp.MustCompile("</?News(:\\w+)>")

func Handler(ctx context.Context) {

	db := s3.New(ctx)

	users, err := model.FindAllUsers(db)
	if err != nil {
		return
	}

	for _, user := range users {
		jobs, e := user.FindAllJobs(db)
		if e != nil {
			continue
		}
		for _, job := range jobs {
			if job.ReadyForWork() {
				Now(db, job)
			}
		}
	}
}

func Now(db s3.Client, job model.Job) {

	if job.Type == model.SitemapJobType {
		_, _ = sitemap.Handle(db, job.UserID, job.URLs[0])
		job.SaveWorkResult(db, nil)
		return
	}

	if job.Type == model.NewsJobType {
		work, err := newsJob(db, job)
		job.SaveWorkResult(db, err)
		if !work.IsEmpty() {
			// todo - use selenium to handle bot defense
		}
		return
	}

	log.Warn().Msgf("unknown job type: %d", job.Type)
}

func newsJob(db s3.Client, job model.Job) (work model.Work, err error) {

	f := func(URL string) (model.Items, error) {

		res, e := http.Get(URL)
		if e != nil {
			return nil, errors.Join(errors.New("failed to http.Get url"), e)
		}
		defer res.Body.Close()

		var b []byte
		if b, e = io.ReadAll(res.Body); e != nil {
			return nil, errors.Join(errors.New("failed to io.ReadAll response body"), e)
		}

		if strings.Contains(URL, "https://www.bing.com/") {
			str := bingRegexp.ReplaceAllStringFunc(string(b), func(s string) string {
				return strings.ReplaceAll(s, ":", "_")
			})
			b = []byte(str)
		}

		var rss struct {
			Channel struct {
				Items model.Items `xml:"item"`
			} `xml:"channel"`
		}
		if e = xml.Unmarshal(b, &rss); e != nil {
			return nil, errors.Join(errors.New("failed to unmarshal xml"), e)
		}

		return rss.Channel.Items, nil
	}

	var keywords []string
	for _, keyword := range job.Keywords {
		keywords = append(keywords, url.QueryEscape(keyword))
	}

	var items model.Items
	for _, u := range job.URLs {
		URL := strings.ReplaceAll(u, "{KEYWORD_QUERY}", strings.Join(keywords, "+"))
		if ii, e := f(URL); e != nil {
			err = errors.Join(err, e)
		} else if ii != nil {
			items = append(items, ii...)
		}
	}

	if len(items) == 0 {
		return
	}

	for _, item := range items {
		if item.IsRedirect() {
			work.Items = append(work.Items, item)
		} else if e := item.Create(db, job); e != nil {
			err = errors.Join(err, e)
		}
	}

	return
}
