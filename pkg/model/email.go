package model

import (
	"encoding/base64"

	"github.com/oklog/ulid/v2"
)

// Email is the starting point for Authenticating a [User].
type Email struct {
	ID     string    `json:"id"`
	UserID ulid.ULID `json:"user_id"`
}

func (e *Email) Dir() string {
	return "email/"
}

func (e *Email) Key() string {
	return e.Dir() + base64.URLEncoding.EncodeToString([]byte(e.ID)) + "_.json"
}
