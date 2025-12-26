package model

import (
	"fmt"

	"github.com/rs/zerolog"
)

type Targets map[string]bool

func (t Targets) MarshalZerologObject(evt *zerolog.Event) {
	var s []string
	for k, v := range t {
		s = append(s, fmt.Sprintf("%s:%v", k, v))
	}
	evt.Strs("targets", s)
}

func (t Targets) String() string {
	var s string
	if t != nil {
		for k, v := range t {
			s += fmt.Sprintf("%s=%v;", k, v)
		}
	}
	return s
}

func (t Targets) None() bool {
	return len(t) == 0
}

func (t Targets) FollowAll() bool {
	if t.None() {
		return false
	}
	_, ok := t["*"]
	return len(t) == 1 && ok
}

func (t Targets) Exist(s string) bool {
	v, k := t[s]
	if t.FollowAll() {
		return !k
	}
	return k && v
}
