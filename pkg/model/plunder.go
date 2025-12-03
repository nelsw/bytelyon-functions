package model

import (
	"bytelyon-functions/pkg/service/em"
	"bytelyon-functions/pkg/service/fn"
	"bytelyon-functions/pkg/service/s3"
	"encoding/json"
	"errors"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Plunder struct {
	*User  `json:"-"`
	ID     ulid.ULID `json:"id"`
	Target string    `json:"target"`
	Follow []string  `json:"follow"`
	Loots  Loots     `json:"loot"`
}

func (p *Plunder) Path() string {
	return p.User.Dir() + "/pw"
}

func (p *Plunder) Dir() string {
	return p.Path() + "/" + p.ID.String()
}

func (p *Plunder) Key() string {
	return p.Dir() + "/_.json"
}

func (p *Plunder) MarshalZerologObject(evt *zerolog.Event) {
	evt.Str("id", p.ID.String()).
		Str("target", p.Target).
		Strs("follow", p.Follow).
		Int("loots", len(p.Loots))
}

func (p *Plunder) Delete() error {

	log.Info().EmbedObject(p).Msg("deleting plunder")

	err := em.Delete(p)
	log.Err(err).EmbedObject(p).Msg("deleted plunder")
	if err != nil {
		return err
	}

	log.Info().EmbedObject(p).Msg("deleting job")

	var j *Job
	if j, err = NewJob(p.User, p.ID).Find(); err != nil {
		log.Warn().Err(err).Msg("failed to find job, it may not exist or have been deleted")
		return nil
	}

	err = em.Delete(j)

	log.Err(err).Msg("deleted job")

	return err
}

func (p *Plunder) Find() error {
	user := p.User
	if err := em.Find(p); err != nil {
		return err
	}
	p.User = user

	keys, _ := s3.New().Keys(p.Dir()+"/loot", "", 1000)

	for _, k := range keys {
		p.Loots = append(p.Loots, NewLoot(p, k))
	}

	return nil
}

func (p *Plunder) Create(b []byte) (*Plunder, error) {

	log.Info().Msgf("creating plunder: %s", string(b))

	var v Plunder
	if err := json.Unmarshal(b, &v); err != nil {
		log.Err(err).Msg("failed to unmarshal plunder")
		return nil, err
	}

	if v.Target == "" {
		return nil, errors.New("missing target")
	}

	v.User = p.User
	v.ID = NewUlid()

	if err := em.Save(&v); err != nil {
		log.Err(err).Msg("failed to save plunder")
		return nil, err
	}

	log.Info().EmbedObject(p).Msg("created plunder")

	return &v, nil
}

func (p *Plunder) FindAll() ([]*Plunder, error) {
	out, err := em.FindAll(p, regexp.MustCompile(`.*/pw/([A-Za-z0-9]{26}/_.json)$`))
	if err != nil {
		return nil, err
	}

	var keys []string
	for idx, val := range out {

		val.User = p.User

		if keys, err = em.Keys(val, regexp.MustCompile(val.Dir())); err != nil {
			log.Warn().Err(err).Msg("failed to find keys")
			continue
		}

		for _, k := range keys {
			if strings.HasSuffix(k, "/_.json") {
				continue
			}
			out[idx].Loots = append(out[idx].Loots, NewLoot(val, k))
		}

		sort.Slice(out[idx].Loots, func(i, j int) bool {
			return out[idx].Loots[i].ID.Timestamp().UnixMilli() > out[idx].Loots[j].ID.Timestamp().UnixMilli()
		})
	}
	return out, nil
}

func (p *Plunder) Work() {

	log.Info().EmbedObject(p).Msg("working plunder")

	if err := p.Find(); err != nil {
		log.Err(err).Msg("failed to find plunder to work")
		return
	}

	log.Trace().EmbedObject(p).Msg("found workable plunder")

	out, err := fn.New().Request("bytelyon-playwrighter", map[string]any{
		"dir":    p.Dir() + "/loot/" + NewUlid().String() + "/",
		"target": p.Target,
		"follow": p.Follow,
	})

	var result string
	if err != nil {
		result = err.Error()
	} else {
		result = string(out)
	}

	log.Info().Str("result", result).Msg("worked plunder")

	var job *Job
	if job, err = NewJob(p.User, p.ID).Find(); err != nil {
		log.Warn().Err(err).Msg("failed to find job")
		return
	}

	job.Results[time.Now().UTC().UnixMilli()] = result
	if err = em.Save(job); err != nil {
		log.Err(err).Msg("failed to save job")
	}

	log.Info().EmbedObject(p).Msg("updated plunder job results")
}

func NewPlunder(user *User, s ...any) *Plunder {

	log.Info().Str("user", user.ID.String()).Msgf("instantiating plunder: %v", s)

	pw := Plunder{User: user}
	if len(s) > 0 {
		if _, ok := s[0].(ulid.ULID); ok {
			pw.ID = s[0].(ulid.ULID)
		} else {
			pw.ID = ulid.MustParse(s[0].(string))
		}
	}

	log.Info().EmbedObject(&pw).Msg("instantiated plunder")

	return &pw
}
