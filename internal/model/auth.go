package model

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"unicode"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
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

func (e *Email) Path() string {
	return fmt.Sprintf("/auth/email/%s", base64.URLEncoding.EncodeToString([]byte(e.ID)))
}

type Password struct {
	UserID ulid.ULID `json:"user_id"`
	Hash   []byte    `json:"hash"`
	Text   string    `json:"-"`
}

func NewPassword(u *User, s string) (*Password, error) {
	p := &Password{UserID: u.ID, Text: s}
	if err := p.Validate(); err != nil {
		return nil, err
	} else if p.Hash, err = bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Password) Path() string {
	return fmt.Sprintf("/auth/pork/%s", p.UserID)
}

func (p *Password) Compare() error {
	return bcrypt.CompareHashAndPassword(p.Hash, []byte(p.Text))
}

func (p *Password) Validate() error {

	if len(p.Text) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	var number, lower, upper, special bool
	for _, c := range p.Text {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			lower = true
		}
	}

	if !lower {
		return errors.New("password must contain at least one lowercase letter")
	} else if !upper {
		return errors.New("password must have at least 1 uppercase character")
	} else if !special {
		return errors.New("password must have at least 1 special character")
	} else if !number {
		return errors.New("password must have at least 1 number character")
	}
	return nil
}
