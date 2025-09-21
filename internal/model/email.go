package model

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"

	"github.com/oklog/ulid/v2"
)

// Email is the starting point for Authenticating a [User].
type Email struct {
	ID       string    `json:"id"`
	UserID   ulid.ULID `json:"user_id"`
	Verified bool      `json:"verified"`
	Token    ulid.ULID `json:"token"`
}

func (e Email) Validate() error {
	if _, err := mail.ParseAddress(e.ID); err != nil {
		return errors.Join(err, errors.New("invalid email address"))
	}
	return nil
}

func (e Email) Path() string {
	return fmt.Sprintf("email/%s", base64.URLEncoding.EncodeToString([]byte(e.ID)))
}

func (e Email) User() User {
	return User{ID: e.UserID}
}
