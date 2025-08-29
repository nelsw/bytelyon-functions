package bot

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var bingRegexp *regexp.Regexp

func init() {
	bingRegexp = regexp.MustCompile("</?News(:\\w+)>")
}

type Work struct {
	ID    ulid.ULID `json:"id"`
	Job   Job       `json:"job"`
	Items Items     `json:"items"`
	Err   error     `json:"error"`
}

func NewWork(j Job) Work {

	w := Work{
		ID:  ulid.Make(),
		Job: j,
	}

	switch j.Type {
	case NewsJobType:
		w.Items, w.Err = workNews(j)
	default:
		w.Err = errors.New(fmt.Sprintf("unknown job type [%d]", j.Type))
	}

	return w
}

func workNews(j Job) (items Items, err error) {

	f := func(u string) (Items, error) {
		u = fmt.Sprintf(u, strings.Join(j.Keywords, ","))
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

		var rss RSS
		if e = xml.Unmarshal(b, &rss); e != nil {
			return nil, errors.Join(errors.New("failed to unmarshal xml"), e)
		}

		return rss.Channel.Items, nil
	}

	for _, u := range j.Type.URLs() {
		ii, e := f(u)
		err = errors.Join(err, e)
		if ii != nil {
			items = append(items, ii...)
		}
		log.Err(e).Any("job", j).Array("items", ii).Send()
	}

	return
}
