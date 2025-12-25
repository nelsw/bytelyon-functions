package internal

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Worker struct {
	stop bool
	done bool
}

func NewWorker() *Worker {
	return new(Worker)
}

func (w *Worker) Start() {

	w.done = false

	log.Info().Msg("worker started")

	users := Users()
	log.Info().Msgf("users found %d", len(users))

	var jobs []*Job
	for _, u := range users {
		jobs = append(jobs, u.FindJobs()...)
	}

	log.Info().Msgf("jobs found %d", len(jobs))
	if len(jobs) == 0 {
		return
	}

	var tasks []*Task
	for _, j := range jobs {
		if j.Type != SearchJobType {
			log.Trace().EmbedObject(j).Msg("not a search job")
			continue
		}
		if !j.Ready() {
			log.Trace().EmbedObject(j).Msg("not ready")
			continue
		}

		s, err := j.FindSearch()
		if err != nil {
			log.Warn().Err(err).EmbedObject(j).Msg("find search failed")
			continue
		}

		tasks = append(tasks, &Task{Job: j, Search: s, Headless: true})
	}

	log.Info().Int("tasks", len(tasks)).Msg("tasks found")

	var wg sync.WaitGroup
	for _, t := range tasks {
		wg.Go(t.execute)
	}

	log.Info().Msg("done jobs")
	wg.Wait()
	log.Info().Msg("worker finished")

	if w.done = true; w.stop {
		return
	}

	time.Sleep(1 * time.Minute)
	w.Start()
}

func (w *Worker) Stop() {
	w.stop = true
}

func (w *Worker) Done() bool {
	return w.stop && w.done
}
