package model

import (
	"github.com/oklog/ulid/v2"
)

type Profile struct {
	UserID ulid.ULID `json:"user_id"`
	Name   string    `json:"name"`
}

func (p *Profile) Dir() string {
	return "user/"
}

func (p *Profile) Key() string {
	return p.Dir() + "profile.json"
}
