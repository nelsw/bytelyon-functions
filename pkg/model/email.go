package model

import (
	"bytelyon-functions/pkg/service/s3"
	"encoding/base64"
	"errors"
	"net/mail"

	"github.com/oklog/ulid/v2"
)

// Email is the starting point for Authenticating a [User].
type Email struct {
	ID     string    `json:"id"`
	UserID ulid.ULID `json:"user_id"`
}

func NewEmail(s string) (*Email, error) {
	if _, err := mail.ParseAddress(s); err != nil {
		return nil, errors.Join(err, errors.New("invalid email address"))
	}
	return &Email{ID: s}, nil
}

func (e *Email) Path() string {
	return "email/" + base64.URLEncoding.EncodeToString([]byte(e.ID))
}

func (e *Email) Key() string {
	return e.Path() + "/_.json"
}

func (e *Email) User() *User {
	return &User{ID: e.UserID}
}

func (e *Email) Find(db s3.Service) error {
	return db.Find(e.Key(), e)
}
