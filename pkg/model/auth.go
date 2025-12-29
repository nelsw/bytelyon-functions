package model

import (
	"bytelyon-functions/pkg/db"
	"encoding/base64"
	"errors"
	"net/mail"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewBasicAuth(s string) (*Auth, error) {

	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(b), ":")
	if len(parts) != 2 {
		return nil, errors.New("invalid basic token; must be base64 encoded '<email>:<password>'")
	}

	return &Auth{
		Username: parts[0],
		Password: parts[1],
	}, nil
}

func (a Auth) Validate() error {
	if _, err := mail.ParseAddress(a.Username); err != nil {
		return errors.Join(err, errors.New("invalid email address"))
	}

	if len(a.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	var number, lower, upper, special bool
	for _, c := range a.Password {
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

func (a Auth) Authenticate() (*User, error) {

	err := a.Validate()
	if err != nil {
		return nil, err
	}

	e := &Email{ID: a.Username}
	if err = db.Find(e); err != nil {
		return nil, err
	}

	u := &User{ID: e.UserID}
	if err = db.Find(u); err != nil {
		return nil, err
	}

	p := &Password{UserID: u.ID}
	if err = db.Find(p); err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword(p.Hash, []byte(a.Password)); err != nil {
		return nil, errors.Join(err, errors.New("invalid password"))
	}

	return u, nil
}
