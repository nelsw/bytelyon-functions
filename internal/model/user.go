package model

import (
	"fmt"

	"github.com/oklog/ulid/v2"
)

type User struct {
	ID ulid.ULID `json:"id"`
}

func (u User) Path() string {
	return fmt.Sprintf("db/user/%s", u.ID)
}

type Profile struct {
	UserId ulid.ULID `json:"user_id"`
	Name   string    `json:"name"`
	Image  string    `json:"image"`
}

func (p Profile) Key() string {
	return fmt.Sprintf("db/user/%s/profile", p.UserId)
}
