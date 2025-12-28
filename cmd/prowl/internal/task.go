package internal

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Task struct {
	Job      *Job
	Search   *Search
	Headless bool
}

func (t *Task) MarshalZerologObject(evt *zerolog.Event) {
	evt.Bool("headless", t.Headless).EmbedObject(t.Search)
}

func (t *Task) execute() {

	p, err := NewProwler(t.Search, t.Headless)
	if err != nil {
		log.Panic().Err(err).EmbedObject(t).Msg(`🦁`)
	}

	if err = p.Google(t.Search); err == nil {
		log.Info().EmbedObject(t).Msg("🦁")
		t.Job.CreateResult()
		return
	}

	if !t.Headless {
		log.Err(err).EmbedObject(t).Msg(`🦁`)
		return
	}

	log.Err(err).EmbedObject(t).Msg(`🦁`)
	t.Headless = false
	t.execute()
}
