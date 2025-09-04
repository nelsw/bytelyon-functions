package model

import (
	"errors"
	"fmt"
	"unicode"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

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
