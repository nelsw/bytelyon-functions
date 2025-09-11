package model

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/s3"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

const JobPath = "job"

type JobType int

const (
	NewsJobType JobType = iota + 1
)

var JobTypes = map[JobType]string{
	NewsJobType: "news",
}

type Jobs []Job
type Job struct {
	ID        ulid.ULID `json:"id"`
	UserID    ulid.ULID `json:"user_id"`
	WorkedAt  time.Time `json:"worked_at"`
	WorkedOk  bool      `json:"worked_ok"`
	Type      JobType   `json:"type"`
	Frequency Frequency `json:"frequency"`
	Name      string    `json:"name"`
	Desc      string    `json:"description"`
	URLs      []string  `json:"urls"`
	Keywords  []string  `json:"keywords"`
}

func (j Job) Path() string {
	return fmt.Sprintf("%s/%s", UserKey(j.UserID), JobPath)
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

func (j Job) ReadyForWork() bool {
	return j.WorkedAt.Add(j.Frequency.Duration()).Before(time.Now().UTC())
}

func (j Job) SaveWorkResult(db s3.Client, err error) {
	j.WorkedAt = time.Now().UTC()
	j.WorkedOk = err == nil
	if e := db.Put(j.Key(), app.MustMarshal(j)); e != nil {
		err = errors.Join(err, e)
	}
	log.Err(err).Any("job", j).Msg("job worked")
	return
}

func (j Job) Items(db s3.Client) (items Items, err error) {
	var keys []string
	fmt.Println(j.Key() + "/" + ItemPath)
	if keys, err = db.Keys(j.Key()+"/"+ItemPath, "", "", 1000); err != nil {
		return
	}
	for _, key := range keys {
		if !strings.Contains(key, ItemPath) {
			continue
		}
		fmt.Println(key)
		var item Item
		if e := db.Find(key, &item); e != nil {
			err = errors.Join(e, e)
		} else {
			items = append(items, item)
		}
	}
	return
}

func JobKey(userID ulid.ULID, jobID any) string {
	return fmt.Sprintf("%s/%s/%s", UserKey(userID), JobPath, jobID)
}
