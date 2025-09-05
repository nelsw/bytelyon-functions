package model

import (
	"fmt"

	"github.com/oklog/ulid/v2"
)

type User struct {
	ID ulid.ULID `json:"id"`
}

func (u User) Path() string {
	return fmt.Sprintf("user")
}

func (u User) Key() string {
	return fmt.Sprintf("%s/%s", u.Path(), u.ID)
}
