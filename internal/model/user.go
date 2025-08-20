package model

import (
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type RoleType int

const (
	Owner RoleType = iota
	Admin
	Basic
)

type User struct {
	ID        ulid.ULID  `json:"id"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	Roles     []RoleType `json:"roles"`
}

type UserProfile struct {
	ID    ulid.ULID `json:"id"` // User.ID
	Name  string    `json:"name"`
	Image string    `json:"image"`
}
type UserEmail struct {
	ID       string    `json:"id" orm:"pk"` // email
	UserID   ulid.ULID `json:"user_id"`
	Verified bool      `json:"verified"`
	Token    string    `json:"token"`
}
type UserPassword struct {
	ID    ulid.ULID `json:"id"` // User.ID
	Value []byte    `json:"value"`
}

func (p *UserPassword) Validate(s string) error {
	return bcrypt.CompareHashAndPassword(p.Value, []byte(s))
}
