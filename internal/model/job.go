package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
)

type JobType int

const (
	NewsJobType JobType = iota + 1
)

var JobTypes = map[JobType]string{
	NewsJobType: "news",
}

type Job struct {
	ID        ulid.ULID `json:"id"`
	UserID    ulid.ULID `json:"user_id"`
	WorkedAt  time.Time `json:"worked_at"`
	Type      JobType   `json:"type"`
	Frequency Frequency `json:"frequency"`
	Name      string    `json:"name"`
	Desc      string    `json:"description"`
	URLs      []string  `json:"urls"`
	Keywords  []string  `json:"keywords"`
}

func (j Job) Path() string {
	return fmt.Sprintf("%s/job", User{j.UserID}.Key())
}

func (j Job) Key() string {
	return fmt.Sprintf("%s/%s", j.Path(), j.ID)
}

func (j Job) Validate() (err error) {

	if _, ok := JobTypes[j.Type]; !ok {
		err = errors.Join(fmt.Errorf("job type must be set"))
	}
	if len(j.Keywords) == 0 {
		err = errors.Join(fmt.Errorf("job keywords must be set"))
	}
	if _, ok := FrequencyUnits[j.Frequency.Unit]; !ok {
		err = errors.Join(fmt.Errorf("job frequency must be one of: [m, h, d]"))
	}
	return
}

func (j Job) Ready() bool {
	return j.WorkedAt.Add(j.Frequency.Duration()).Before(time.Now().UTC())
}

func (j Job) MarshalZerologObject(e *zerolog.Event) {
	e.Stringer("id", j.ID).
		Str("type", JobTypes[j.Type])
}
