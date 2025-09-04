package model

import (
	"fmt"

	"github.com/oklog/ulid/v2"
)

type User struct {
	ID ulid.ULID `json:"id"`
}

func (u User) Path() string {
	return fmt.Sprintf("user/%s", u.ID)
}
