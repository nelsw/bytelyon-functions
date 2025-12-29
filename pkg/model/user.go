package model

import (
	"github.com/oklog/ulid/v2"
)

type User struct {
	ID ulid.ULID `json:"id"`
}

func (u *User) Dir() string {
	return "user/"
}

func (u *User) Key() string {
	return u.Dir() + "_.json"
}
