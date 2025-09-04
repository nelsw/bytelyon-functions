package model

import (
	"fmt"

	"github.com/oklog/ulid/v2"
)

type User struct {
	ID ulid.ULID `json:"id"`
}

func (u User) Path() string {
	return fmt.Sprintf("/user/%s/user", u.ID)
}

func NewUser() *User {
	return &User{NewUlid()}
}

type Profile struct {
	UserID ulid.ULID `json:"user_id"`
	Name   string    `json:"name"`
	Image  string    `json:"image"`
}

func (p Profile) Key() string {
	return fmt.Sprintf("/user/%s/profile", p.UserID)
}

func NewProfile(u *User) *Profile {
	return &Profile{UserID: u.ID}
}
