package model

import (
	"bytelyon-functions/internal/util"
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

func NewEmail(u *User, s string) (*Email, error) {
	e := &Email{
		UserID: u.ID,
		ID:     s,
		Token:  NewUlid(),
	}
	if err := e.Validate(); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *Email) Validate() error {
	if _, err := mail.ParseAddress(e.ID); err != nil {
		return errors.Join(err, errors.New("invalid email address"))
	}
	return nil
}

func (e *Email) Key() string {
	return fmt.Sprintf("%s/db/auth/email/%s.json", util.AppMode(), base64.URLEncoding.EncodeToString([]byte(e.ID)))
}
