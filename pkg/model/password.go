package model

import (
	"github.com/oklog/ulid/v2"
)

type Password struct {
	UserID ulid.ULID `json:"user_id"`
	Hash   []byte    `json:"hash"`
}

func (p *Password) Dir() string {
	return "user/"
}

func (p *Password) Key() string {
	return p.Dir() + "pork.json"
}
