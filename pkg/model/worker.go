package model

import (
	"bytelyon-functions/pkg/db"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Worker struct {
	stop, done bool
}

func (w *Worker) Start() {

	w.done = false

	users, _ := db.List(&User{})

	log.Info().
		Int("count", len(users)).
		Msg("Worker - Found Users")

	var prowlers []*Prowler
	prowlers = append(prowlers, w.list(users, SearchProwlType)...)
	prowlers = append(prowlers, w.list(users, SitemapProwlType)...)
	prowlers = append(prowlers, w.list(users, ArticleProwlType)...)

	log.Info().
		Int("count", len(prowlers)).
		Msg("Worker - Found Prowlers")

	var wg sync.WaitGroup
	for _, prowler := range prowlers {
		wg.Go(prowler.Prowl)
	}
	wg.Wait()

	if w.done = true; w.stop {
		return
	}

	time.Sleep(time.Minute)

	w.Start()
}

func (w *Worker) Stop() {
	w.stop = true
}

func (w *Worker) Done() bool {
	return w.stop && w.done
}

func (w *Worker) list(users []*User, t ProwlerType) []*Prowler {

	var wg sync.WaitGroup

	var prowlers []*Prowler

	for _, user := range users {

		wg.Go(func() {

			p, _ := db.List(&Prowler{UserID: user.ID, Type: t})

			log.Trace().
				Stringer("user", user.ID).
				Int("count", len(p)).
				Msgf("Worker - Found [%s] Prowlers", t)

			prowlers = append(prowlers, p...)
		})
	}

	wg.Wait()

	log.Debug().
		Int("count", len(prowlers)).
		Msgf("Worker - Found [%s] Prowlers", t)

	return prowlers
}
