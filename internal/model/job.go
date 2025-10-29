package model

import (
	"bytelyon-functions/internal/client/s3"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/oklog/ulid/v2"
)

type JobType int

const (
	NewsJobType JobType = iota + 1
	SitemapJobType
)

var JobTypes = map[JobType]string{
	NewsJobType:    "news",
	SitemapJobType: "sitemap",
}

type Jobs []Job
type Job struct {
	User      *User     `json:"-"`
	ID        ulid.ULID `json:"id"`
	WorkedAt  time.Time `json:"worked_at"`
	Type      JobType   `json:"type"`
	Frequency Frequency `json:"frequency"`
	Name      string    `json:"name"`
	Desc      string    `json:"description"`
	URLs      []string  `json:"urls"`
	Keywords  []string  `json:"keywords"`
	Items     Items     `json:"items"`
}

func NewJob(req events.APIGatewayV2HTTPRequest) (*Job, error) {

	u, err := NewUser(req)
	if err != nil {
		return nil, err
	}

	var j Job
	if req.Body != "" {
		if err = json.Unmarshal([]byte(req.Body), &j); err != nil {
			return nil, err
		}
	}

	j.User = u

	if j.ID.IsZero() && req.QueryStringParameters["id"] != "" {
		j.ID = ulid.MustParse(req.QueryStringParameters["id"])
	}

	return &j, err
}

func (j *Job) Validate() (err error) {

	if _, ok := JobTypes[j.Type]; !ok {
		err = errors.Join(fmt.Errorf("job type must be set"))
	}
	if j.Type == SitemapJobType && len(j.URLs) == 0 {
		err = errors.Join(fmt.Errorf("job url must be set"))
	}
	if len(j.Keywords) == 0 {
		err = errors.Join(fmt.Errorf("job keywords must be set"))
	}
	if _, ok := FrequencyUnits[j.Frequency.Unit]; !ok {
		err = errors.Join(fmt.Errorf("job frequency must be one of: [m, h, d]"))
	}
	return
}

func (j *Job) Path() string {
	return j.User.Path() + "/job"
}

func (j *Job) Key() string {
	return j.Path() + "/" + j.ID.String() + "/_.json"
}

func (j *Job) Create() (*Job, error) {
	j.ID = NewUlid()
	return j.save(true)
}

func (j *Job) Update() (*Job, error) {
	return j.save(false)
}

func (j *Job) Delete() (any, error) {
	return nil, s3.New().Delete(j.Key())
}

func (j *Job) FindAll() (Jobs, error) {
	db := s3.New()

	keys, err := db.Keys(j.Path(), "", 1000)
	if err != nil {
		return nil, err
	}

	var vv Jobs

	for _, k := range keys {

		o, e := db.Get(k)
		if e != nil {
			err = errors.Join(err, e)
			continue
		}

		var v Job
		if e = json.Unmarshal(o, &v); e != nil {
			err = errors.Join(err, e)
			continue
		}

		vv = append(vv, v)
	}

	return vv, err
}

func (j *Job) save(run bool) (*Job, error) {

	if err := j.Validate(); err != nil {
		return nil, err
	}

	if j.Type == NewsJobType {
		var keywordQuery string
		for i, keyword := range j.Keywords {
			if i > 0 {
				keywordQuery += "+"
			}
			keywordQuery += url.QueryEscape(keyword)
		}
		j.URLs = []string{
			fmt.Sprintf("https://news.google.com/rss/search?q=%s&hl=en-US&gl=US&ceid=US:en", keywordQuery),
			fmt.Sprintf("https://www.bing.com/news/search?format=rss&q=%s", keywordQuery),
			fmt.Sprintf("https://www.bing.com/search?format=rss&q=%s", keywordQuery),
		}
	}

	if run {

		if j.Type == NewsJobType {
			j.doNewsWork()
			j.WorkedAt = time.Now().UTC()
		}
	}

	if b, err := json.Marshal(j); err != nil {
		return nil, err
	} else if err = s3.New().Put(j.Key(), b); err != nil {
		return nil, err
	}

	return j, nil
}

func (j *Job) doNewsWork() (err error) {

	var items Items

	for _, u := range j.URLs {

		rss, e := NewRSS(u)
		if e != nil {
			err = errors.Join(err, e)
		}

		if rss.Channel.Items == nil {
			continue
		}

		for _, i := range rss.Channel.Items {
			i.Scrub()
			items = append(items, i)
		}
	}

	j.Items = append(items, j.Items...)

	return
}
