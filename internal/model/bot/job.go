package bot

import (
	"bytelyon-functions/internal/entity"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type JobType int

const (
	Custom JobType = iota + 1
	GoogleNews
	BingNews
)

func (v JobType) URLs() []string {
	switch v {
	case GoogleNews:
		return []string{"https://news.google.com/rss/search?q=%s&hl=en-US&gl=US&ceid=US:en"}
	case BingNews:
		return []string{"https://www.bing.com/search?format=rss&q=", "https://www.bing.com/news/search?format=rss&q="}
	default:
		return nil
	}
}

func (v JobType) Sanitize(b []byte) []byte {

	if v != BingNews {
		return b
	}

	reg := regexp.MustCompile("</?News(:\\w+)>")
	str := reg.ReplaceAllStringFunc(string(b), func(s string) string {
		return strings.ReplaceAll(s, ":", "_")
	})

	return []byte(str)
}

type Job struct {
	ID        ulid.ULID `json:"id"`
	Name      string    `json:"name"`
	Type      JobType   `json:"type"`
	Keywords  []string  `json:"keywords"`
	Frequency Frequency `json:"frequency"`
	Error     error     `json:"error"`
	LastRunAt Timestamp `json:"last_run_at"`
	CreatedAt Timestamp `json:"created_at"`
	UpdatedAt Timestamp `json:"updated_at"`
}

type Timestamp struct {
	UTC   time.Time `json:"utc"`
	Unix  int64     `json:"unix"`
	Human string    `json:"human"`
}

func NewTimeStamp() Timestamp {
	t := time.Now().UTC()
	return Timestamp{
		UTC:   t,
		Unix:  t.Unix(),
		Human: t.Format("01/02/06 3:04PM"),
	}
}

func (t Timestamp) BeforeNow(d time.Duration) bool {
	return t.UTC.Add(d).Before(time.Now().UTC())
}

type FrequencyUnit string

const (
	Minute FrequencyUnit = "m"
	Hour   FrequencyUnit = "h"
	Day    FrequencyUnit = "d"
)

var FrequencyUnits = map[string]FrequencyUnit{
	string(Minute): Minute,
	string(Hour):   Hour,
	string(Day):    Day,
}

type Frequency struct {
	Unit  FrequencyUnit `json:"unit"`
	Value int           `json:"value"`
}

func (f Frequency) Duration() time.Duration {
	if f.Unit == Minute {
		return time.Minute * time.Duration(f.Value)
	} else if f.Unit == Hour {
		return time.Hour * time.Duration(f.Value)
	} else { // day
		return time.Hour * 24 * time.Duration(f.Value)
	}
}

func (v *Job) Validate() error {
	if v.Type == 0 {
		return fmt.Errorf("job type must be set")
	} else if len(v.Keywords) == 0 {
		return fmt.Errorf("job keywords must be set")
	} else if _, ok := FrequencyUnits[string(v.Frequency.Unit)]; !ok {
		return fmt.Errorf("job frequency must be one of: [m, h, d]")
	}
	return nil
}

func (v *Job) Ready() bool {
	return v.LastRunAt.BeforeNow(v.Frequency.Duration())
}

func (v *Job) CreateWork() error {

	log.Info().Str("ID", v.ID.String()).Msg("creating work")
	if v.Type == Custom {
		v.Error = errors.New("custom job not supported yet")
	} else {
		for _, u := range v.Type.URLs() {
			if v.Error = createWorkItems(v.ID, v.Type, u, v.Keywords); v.Error != nil {
				break
			}
		}
	}

	v.UpdatedAt = NewTimeStamp()
	if err := entity.New().Value(v).Save(); err != nil {
		log.Error().Err(err).Str("ID", v.ID.String()).Msg("failed to update job last run time")
		return err
	}

	log.Info().Str("ID", v.ID.String()).Msg("job updated")

	return nil
}
