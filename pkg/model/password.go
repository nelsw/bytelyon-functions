package model

import (
	"bytelyon-functions/pkg/service/s3"
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	User *User  `json:"-"`
	Text string `json:"-"`
	Hash []byte `json:"hash"`
}

func NewPassword(s string) (*Password, error) {

	if len(s) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	var number, lower, upper, special bool
	for _, c := range s {
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
		return nil, errors.New("password must contain at least one lowercase letter")
	} else if !upper {
		return nil, errors.New("password must have at least 1 uppercase character")
	} else if !special {
		return nil, errors.New("password must have at least 1 special character")
	} else if !number {
		return nil, errors.New("password must have at least 1 number character")
	}

	return &Password{
		Text: s,
	}, nil
}

func (p *Password) Compare(b []byte) error {
	return bcrypt.CompareHashAndPassword(b, []byte(p.Text))
}

func (p *Password) Path() string {
	return p.User.Dir() + "/pork"
}

func (p *Password) Key() string {
	return p.Path() + "/_.json"
}

func (p *Password) Find(db s3.Service, u *User) error {

	p.User = u

	var pass Password
	if err := db.Find(p.Key(), &pass); err != nil {
		return err
	} else if err = p.Compare(pass.Hash); err != nil {
		return err
	}

	return nil
}
