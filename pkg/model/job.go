package model

import (
	"bytelyon-functions/pkg/service/em"
	"encoding/json"
	"errors"
	"maps"
	"regexp"
	"slices"
	"sort"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type JobType string

const (
	NewsJobType    JobType = "news"
	SearchJobType  JobType = "search"
	SitemapJobType JobType = "sitemap"
	PlunderJobType JobType = "plunder"
)

var (
	validJobTypes = regexp.MustCompile(`^(news|search|sitemap|plunder)$`)
)

type Job struct {
	User    *User         `json:"-"`
	ID      ulid.ULID     `json:"id"`
	Name    string        `json:"name"`
	Type    JobType       `json:"type"`
	Freq    *Frequency    `json:"frequency"`
	Results map[int64]any `json:"results"`
}

func NewJob(args ...any) *Job {
	var v = new(Job)
	for _, arg := range args {
		switch arg.(type) {
		case *User:
			v.User = arg.(*User)
		case ulid.ULID:
			v.ID = arg.(ulid.ULID)
		case string:
			if id, err := ulid.Parse(arg.(string)); err == nil {
				v.ID = id
			}
		}
	}
	return v
}

func (j *Job) Path() string {
	return j.User.Dir() + "/job"
}

func (j *Job) Dir() string {
	return j.Path() + "/" + j.ID.String()
}

func (j *Job) Key() string {
	return j.Dir() + "/_.json"
}

func (j *Job) Save(b []byte) (*Job, error) {

	var v Job
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, err
	}
	v.User = j.User

	if !validJobTypes.MatchString(string(v.Type)) {
		return nil, errors.New("invalid job type, must be one of: news, search, sitemap, plunder")
	} else if v.Freq == nil {
		return nil, errors.New("missing frequency")
	} else if err := v.Freq.Validate(); err != nil {
		return nil, err
	} else if v.ID.IsZero() {
		return nil, errors.New("missing job id")
	} else if w := v.Worker(); w == nil {
		return nil, errors.New("invalid worker id/type combination")
	}

	// todo - validate the worker exists

	// load existing results
	if f, err := v.Find(); err == nil {
		v.Results = f.Results
	}

	if err := em.Save(&v); err != nil {
		return nil, err
	}
	v.User = j.User
	return &v, nil
}

func (j *Job) Delete() error {
	return em.Delete(j)
}

func (j *Job) FindAll() ([]*Job, error) {
	jobs, err := em.FindAll(j, regexp.MustCompile(`.*/job/([A-Za-z0-9]{26}/_.json)$`))
	if err != nil {
		log.Warn().Err(err).Msg("failed to find jobs")
		return []*Job{}, err
	}
	for i, v := range jobs {
		v.User = j.User
		jobs[i] = v
	}
	return jobs, nil
}

func (j *Job) Find() (*Job, error) {

	v := NewJob(j.User, j.ID)
	if err := em.Find(v); err != nil {
		return nil, err
	}

	v.User = j.User

	return v, nil
}

func (j *Job) Ready() bool {
	if len(j.Results) == 0 {
		return true
	}
	keys := slices.Collect(maps.Keys(j.Results))
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	return time.UnixMilli(keys[0]).Add(j.Freq.Duration()).After(time.Now().UTC())
}

func (j *Job) Worker() Worker {

	switch j.Type {
	case NewsJobType:
		return NewNews(j.User, j.ID)
	case PlunderJobType:
		return NewPlunder(j.User, j.ID)
	case SitemapJobType:
		log.Warn().Msg("sitemap worker not yet implemented")
	case SearchJobType:
		log.Warn().Msg("search worker not yet implemented")
	default:
		log.Warn().Msgf("invalid job type: %s", j.Type)
	}

	return nil
}
