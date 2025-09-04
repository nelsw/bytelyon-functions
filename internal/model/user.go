package model

import (
	"bytelyon-functions/internal/util"
	"fmt"

	"github.com/oklog/ulid/v2"
)

type User struct {
	ID ulid.ULID `json:"id"`
}

func (u User) Path() string {
	return fmt.Sprintf("%s/db/user/%s", util.AppMode(), u.ID)
}

func (u User) Key() string {
	return fmt.Sprintf("%s.json", u.Path())
}

type Profile struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func (p Profile) Key(user User) string {
	return fmt.Sprintf("%s_profile.json", user.Path())
}
