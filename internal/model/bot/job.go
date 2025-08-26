package bot

import (
	"bytelyon-functions/internal/entity"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type JobType int

const (
	GoogleNews JobType = iota + 1
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
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      JobType   `json:"type"`
	Roots     []string  `json:"roots"`
	Keywords  []string  `json:"keywords"`
	Frequency int       `json:"frequency"` // minutes
	Success   bool      `json:"success"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (v *Job) Validate() error {
	if v.Type == 0 {
		return fmt.Errorf("job type must be set")
	} else if len(v.Keywords) == 0 {
		return fmt.Errorf("job keywords must be set")
	} else if v.Frequency == 0 {
		return fmt.Errorf("job frequency must be set")
	}
	return nil
}

func (v *Job) Ready() bool {
	return v.UpdatedAt.Add(time.Minute * time.Duration(v.Frequency)).Before(time.Now())
}

func (v *Job) CreateWork() error {

	log.Info().Str("ID", v.ID.String()).Msg("creating work")

	ok := true
	for _, u := range v.Type.URLs() {
		if err := createWorkItems(v.ID, v.Type, u, v.Keywords); err != nil {
			ok = false
			break
		}
	}

	v.Success = ok
	v.UpdatedAt = time.Now()
	if err := entity.New().Value(v).Save(); err != nil {
		log.Error().Err(err).Str("ID", v.ID.String()).Msg("failed to update job last run time")
		return err
	}

	log.Info().Str("ID", v.ID.String()).Msg("job updated")

	return nil
}
