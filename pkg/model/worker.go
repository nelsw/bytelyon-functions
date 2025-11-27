package model

import (
	"sync"

	"github.com/rs/zerolog/log"
)

type Worker interface {
	Work()
}

func DoWork() {
	users, _ := FindAllUsers()

	var workers []Worker
	for _, u := range users {
		jobs, _ := NewJob(u).FindAll()
		for _, j := range jobs {
			if j.Ready() {
				workers = append(workers, j.Worker())
			}
		}
	}

	if len(workers) == 0 {
		log.Info().Msg("no jobs ready, will try again later")
		return
	}

	log.Info().Int("size", len(workers)).Msg("jobs ready")

	var wg sync.WaitGroup

	for _, w := range workers {
		wg.Go(w.Work)
	}

	log.Info().Msg("waiting on work to complete")

	wg.Wait()

	log.Info().Msg("work completed")
}
