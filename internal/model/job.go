package model

import (
	"bytelyon-functions/internal/app"
	"encoding/json"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
)

type Job struct {
	User      User      `json:"-" fake:"skip"`
	Work      []Work    `json:",omitempty"`
	ID        ulid.ULID `json:"id" fake:"skip"`
	WorkID    ulid.ULID `json:"work_id" fake:"skip"`
	Type      JobType   `json:"type"`
	Frequency Frequency `json:"frequency"`
	Name      string    `json:"name"`
	Desc      string    `json:"description"`
	URLs      []string  `json:"urls"`
	Keywords  []string  `json:"keywords"`
}

func (j Job) Path() string {
	return fmt.Sprintf("%s/job", j.User.Key())
}

func (j Job) Key() string {
	return fmt.Sprintf("%s/%s", j.Path(), j.ID)
}

func (j Job) Validate() error {
	if _, ok := JobTypes[j.Type]; !ok {
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
		Str("type", JobTypes[j.Type])
}

func MakeJob(u User, b []byte) (j Job, err error) {
	if err = json.Unmarshal(b, &j); err == nil {
		err = j.Validate()
	}
	if err == nil {
		j.User = u
		if j.ID.IsZero() {
			j.ID = app.NewUlid()
		}
	}
	return
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
