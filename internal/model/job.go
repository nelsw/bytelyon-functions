package model

import (
	"bytelyon-functions/pkg/entity"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Job struct {
	ID        ulid.ULID `json:"id" fake:"skip"`
	UserID    ulid.ULID `json:"user_id" fake:"skip"`
	WorkID    ulid.ULID `json:"work_id" fake:"skip"`
	Type      JobType   `json:"type"`
	Frequency Frequency `json:"frequency"`
	Err       error     `json:"error" fake:"skip"`
	Name      string    `json:"name"`
	Desc      string    `json:"description"`
	URLs      []string  `json:"urls"`
	Keywords  []string  `json:"keywords"`
	User      User      `json:"-"`
}

func (j Job) Path() string {
	return fmt.Sprintf("%s/job/%s", j.User.Path(), j.ID)
}

func (j Job) Validate() error {
	if j.ID.IsZero() {
		return fmt.Errorf("job id is not set")
	} else if j.UserID.IsZero() {
		return fmt.Errorf("user id is not set")
	} else if _, ok := JobTypes[j.Type]; !ok {
		return fmt.Errorf("job type must be set")
	} else if len(j.Keywords) == 0 {
		return fmt.Errorf("job keywords must be set")
	} else if _, ok = FrequencyUnits[j.Frequency.Unit]; !ok {
		return fmt.Errorf("job frequency must be one of: [m, h, d]")
	}
	return nil
}

func (j Job) Ready() bool {
	now := time.Now().UTC()
	dur := j.Frequency.Duration()
	return j.WorkID.Timestamp().Add(dur).Before(now)
}

func (j Job) MarshalZerologObject(e *zerolog.Event) {
	e.Stringer("id", j.ID).
		Str("type", JobTypes[j.Type]).
		Err(j.Err)
}

func NewJob(userID ulid.ULID) Job {
	return Job{
		ID:     NewUlid(),
		UserID: userID,
	}
}

type JobType int

const (
	NewsJobType JobType = iota + 1
)

var JobTypes = map[JobType]string{
	NewsJobType: "news",
}

func (v JobType) URLs() []string {
	switch v {
	case NewsJobType:
		return []string{
			"https://news.google.com/rss/search?q=%s&hl=en-US&gl=US&ceid=US:en",
			"https://www.bing.com/news/search?format=rss&q=%s",
			"https://www.bing.com/search?format=rss&q=%s",
		}
	default:
		return nil
	}
}

type FrequencyUnit string

const (
	Minute FrequencyUnit = "m"
	Hour   FrequencyUnit = "h"
	Day    FrequencyUnit = "d"
)

var FrequencyUnits = map[FrequencyUnit]time.Duration{
	Minute: time.Minute,
	Hour:   time.Hour,
	Day:    time.Hour * 24,
}

type Frequency struct {
	Unit  FrequencyUnit `json:"unit"`
	Value int           `json:"value"`
}

func (f Frequency) Duration() time.Duration {
	return FrequencyUnits[f.Unit] * time.Duration(f.Value)
}

func (j Job) CreateWork() error {

	w := NewWork(j)
	j.WorkID = w.ID
	j.Err = entity.New().Value(&w).Save()

	log.Err(j.Err).Any("job", j).Msg("worked")

	return entity.New().Value(&j).Save()
}
