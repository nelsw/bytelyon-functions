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

func FindJobs(db s3.Client, userID ulid.ULID) (page Page, err error) {

	var keys []string
	if keys, err = db.Keys(UserKey(userID)+"/job", "", 1000); err != nil {
		return
	}
	page.Total = len(keys)

	var job Job
	var items Items
	for _, key := range keys {
		if !strings.Contains(key, "/item/") {
			if len(items) > 0 {
				page = page.Add(map[string]any{
					"job":   job,
					"items": items,
				})
				items = nil
				job = Job{}
			}

			_ = db.Find(key, &job)
			continue
		}

		var item Item
		_ = db.Find(key, &item)
		items = append(items, item)
	}
	if len(items) > 0 {
		page = page.Add(map[string]any{
			"job":   job,
			"items": items,
		})
	}
	log.Err(err).
		Int("total", page.Total).
		Int("size", page.Size).
		Msg("FindJobs")
	return page, err
}

type Job struct {
	ID        ulid.ULID `json:"id"`
	WorkedAt  time.Time `json:"worked_at,omitempty"`
	WorkedOk  bool      `json:"worked_ok"`
	Type      JobType   `json:"type"`
	Frequency Frequency `json:"frequency"`
	Name      string    `json:"name"`
	Desc      string    `json:"description"`
	URLs      []string  `json:"urls"`
	Keywords  []string  `json:"keywords"`
}

func (j Job) Validate() (err error) {

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

func (j Job) ReadyForWork() bool {
	return j.WorkedAt.Add(j.Frequency.Duration()).Before(time.Now().UTC())
}

func (j Job) SaveWorkResult(db s3.Client, err error, userID ulid.ULID) {

	j.WorkedAt = time.Now().UTC()
	j.WorkedOk = err == nil

	if e := db.Put(JobKey(userID, j.ID), app.MustMarshal(j)); e != nil {
		err = errors.Join(err, e)
	}

	log.Err(err).Any("job", j).Msg("job worked")

	return
}

func (j Job) DoWork(db s3.Client, userID ulid.ULID) {
	if j.Type == SitemapJobType {
		_, err := NewSitemap(j.URLs[0]).Create(db, userID)
		j.SaveWorkResult(db, err, userID)
		return
	}

	if j.Type == NewsJobType {
		err := j.newsJob(db, userID)
		j.SaveWorkResult(db, err, userID)

		return
	}

	log.Warn().Msgf("unknown job type: %d", j.Type)
}

func (j Job) newsJob(db s3.Client, userID ulid.ULID) (err error) {
	var items Items
	for _, u := range j.URLs {

		if ii, e := MakeNewsItems(u); e != nil {
			err = errors.Join(err, e)
		} else if ii != nil {
			items = append(items, ii...)
		}
	}

	for _, item := range items {
		if e := item.Create(db, JobKey(userID, j.ID)); e != nil {
			err = errors.Join(err, e)
		}
	}

	return
}

func JobPath(userID ulid.ULID) string {
	return fmt.Sprintf("%s/%s", UserKey(userID), "job")
}

func JobKey(userID ulid.ULID, jobID any) string {
	return fmt.Sprintf("%s/%s", JobPath(userID), jobID)
}
