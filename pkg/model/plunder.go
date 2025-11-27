package model

import (
	"bytelyon-functions/pkg/service/em"
	"bytelyon-functions/pkg/service/fn"
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
		Array("loots", p.Loots)
}

func (p *Plunder) Delete() error {
	if err := em.Delete(p); err != nil {
		return err
	}

	// delete the associated job
	j, err := NewJob(p.User, p.ID).Find()
	if err != nil {
		log.Warn().Err(err).Msg("failed to find job, it may not exist or have been deleted")
		return nil
	}

	return em.Delete(j)
}

func (p *Plunder) Find() error {
	user := p.User
	if err := em.Find(p); err != nil {
		return err
	}
	p.User = user
	return nil
}

func (p *Plunder) Create(b []byte) (*Plunder, error) {

	var v Plunder
	if err := json.Unmarshal(b, &v); err != nil {
		log.Err(err).Msg("failed to unmarshal p")
		return nil, err
	}

	if v.Target == "" {
		return nil, errors.New("missing target")
	}

	v.User = p.User
	v.ID = NewUlid()

	if err := em.Save(&v); err != nil {
		log.Err(err).Msg("failed to save p")
		return nil, err
	}

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
			out[idx].Loots = append(out[idx].Loots, NewLoot(k))
		}

		sort.Slice(out[idx].Loots, func(i, j int) bool {
			return out[idx].Loots[i].Time > out[idx].Loots[j].Time
		})
	}
	return out, nil
}

func (p *Plunder) Work() {

	if err := p.Find(); err != nil {
		log.Err(err).Msg("failed to find news")
		return
	}

	job, err := NewJob(p.User, p.ID).Find()
	if err != nil {
		log.Err(err).Msg("failed to find job")
		return
	}

	var out []byte
	out, err = fn.New().Request("bytelyon-function-playwrighter", map[string]any{
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

	job.Results[time.Now().UTC().UnixMilli()] = result
	if err = em.Save(job); err != nil {
		log.Err(err).Msg("failed to save job")
	}
}

func NewPlunder(user *User, s ...any) *Plunder {
	pw := Plunder{User: user}
	if len(s) > 0 {
		if _, ok := s[0].(ulid.ULID); ok {
			pw.ID = s[0].(ulid.ULID)
		} else {
			pw.ID = ulid.MustParse(s[0].(string))
		}
	}
	return &pw
}
