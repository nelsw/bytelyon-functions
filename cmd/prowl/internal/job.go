package internal

import (
	"maps"
	"slices"
	"sort"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type JobType string

const (
	SearchJobType JobType = "search"
)

type Job struct {
	UserID  ulid.ULID     `json:"user_id"`
	ID      ulid.ULID     `json:"id"`
	Type    JobType       `json:"type"`
	Freq    *Frequency    `json:"frequency"`
	Results map[int64]any `json:"results"`
}

func (j *Job) String() string {
	return Path("user", j.UserID, "job", j.ID)
}

func (j *Job) MarshalZerologObject(evt *zerolog.Event) {
	evt.Stringer("user", j.UserID).
		Stringer("job", j.ID).
		Str("type", string(j.Type)).
		Stringer("freq", j.Freq).
		Int("results", len(j.Results))
}

func (j *Job) Ready() bool {
	if len(j.Results) == 0 {
		return true
	}
	if j.Freq == nil {
		return false
	}
	keys := slices.Collect(maps.Keys(j.Results))
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	return keys[0]+j.Freq.Duration().Milliseconds() < time.Now().UTC().UnixMilli()
}

func (j *Job) FindSearch() (*Search, error) {
	s := &Search{UserID: j.UserID, ID: j.ID}
	log.Trace().EmbedObject(s).Msg("find search")

	if err := S3.Find(s); err != nil {
		log.Err(err).EmbedObject(s).Msg("find search failed")
		return nil, err
	}

	log.Trace().EmbedObject(s).Msg("found search")
	return s, nil
}

func (j *Job) CreateResult(e ...error) {
	log.Trace().EmbedObject(j).Bool("ok", e == nil).Msg("create job result")
	if k := time.Now().UTC().UnixMilli(); len(e) > 0 {
		j.Results[k] = e[0]
	} else {
		j.Results[k] = nil
	}

	if err := S3.Put(j); err != nil {
		log.Err(err).EmbedObject(j).Msg("failed to create job result")
		return
	}
	log.Info().EmbedObject(j).Bool("ok", e == nil).Msg("created job result")
}

//func (j *Job) Work() {
//
//	if j.Type != SearchJobType {
//		log.Panic().EmbedObject(j).Msg("work called on non search job")
//		return
//	}
//
//	s, err := j.FindSearch()
//	if err != nil {
//		log.Panic().Err(err).EmbedObject(j).Msg("work failed")
//		return
//	}
//
//	ts := time.Now().UTC()
//	err = s.Execute(ts, true)
//	if err == nil {
//
//		return
//	}
//
//	err = s.Execute(ts, false)
//	if err != nil {
//		log.Panic().Err(err).EmbedObject(j).Msg("work failed")
//	}
//}
