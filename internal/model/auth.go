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
	Token    ulid.ULID `json:"token"`
	Verified bool      `json:"verified"`
}

func (e Email) Validate() error {
	if _, err := mail.ParseAddress(e.ID); err != nil {
		return errors.Join(err, errors.New("invalid email address"))
	}
	return nil
}

func (e Email) Path() string {
	return fmt.Sprintf("db/auth/%s", base64.URLEncoding.EncodeToString([]byte(e.ID)))
}

type PasswordText string
type PasswordHash []byte

func (p PasswordText) Validate() error {

	var number, lower, upper, special bool
	for _, c := range p {
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

	if len(p) < 8 {
		return errors.New("password must be at least 8 characters")
	} else if !lower {
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

func (h PasswordHash) Compare(t PasswordText) error {
	return bcrypt.CompareHashAndPassword(h, []byte(t))
}
