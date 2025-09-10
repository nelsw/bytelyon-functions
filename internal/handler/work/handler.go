package work

import (
	"bytelyon-functions/internal/model"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var bingRegexp = regexp.MustCompile("</?News(:\\w+)>")

func Handler() {
	// todo - get all jobs and work all that are ready
}

func Job(job model.Job) error {
	switch job.Type {
	case model.NewsJobType:
		return newsJob(job)
	default:
		return fmt.Errorf("unknown job type [%d]", job.Type)
	}
}

func newsJob(job model.Job) (err error) {

	f := func(u string) (model.Items, error) {

		// todo - handle spaces
		u = fmt.Sprintf(u, strings.Join(job.Keywords, ","))

		res, e := http.Get(u)
		if e != nil {
			return nil, errors.Join(errors.New("failed to http.Get url"), e)
		}
		defer res.Body.Close()

		var b []byte
		if b, e = io.ReadAll(res.Body); e != nil {
			return nil, errors.Join(errors.New("failed to io.ReadAll response body"), e)
		}

		if strings.Contains(u, "https://www.bing.com/") {
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

	var items model.Items
	for _, u := range job.URLs {
		ii, e := f(u)
		if e != nil {
			err = errors.Join(err, e)
		}
		if ii != nil {
			items = append(items, ii...)
		}
	}

	if len(items) == 0 {
		return err
	}

	// todo - put a work object on the bus
	return
}
