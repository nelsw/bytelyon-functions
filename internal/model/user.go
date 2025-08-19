package model

import (
	"encoding/base64"
	"errors"
	"net/mail"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type Access int

const (
	Login Access = iota
	Signup
)

type Credentials struct {
	ID            ulid.ULID `json:"-"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	*User         `json:"-"`
	*UserEmail    `json:"-"`
	*UserPassword `json:"-"`
}

func NewCredentials(token string) (*Credentials, error) {
	if strings.HasPrefix(token, "Basic ") {
		token = strings.TrimPrefix(token, "Basic ")
	}

	b, _ := base64.StdEncoding.DecodeString(token)
	parts := strings.Split(string(b), ":")
	if len(parts) != 2 {
		return nil, errors.New("invalid token")
	}

	// validate email
	if _, err := mail.ParseAddress(parts[0]); err != nil {
		return nil, errors.New("invalid email address")
	}

	// validate password
	var number, lower, upper, special bool
	for _, c := range parts[1] {
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

	if len(parts[1]) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	} else if !lower {
		return nil, errors.New("password must contain at least one lowercase letter")
	} else if !upper {
		return nil, errors.New("password must have at least 1 uppercase character")
	} else if !special {
		return nil, errors.New("password must have at least 1 special character")
	} else if !number {
		return nil, errors.New("password must have at least 1 number character")
	}

	return &Credentials{
		Email:    parts[0],
		Password: parts[1],
	}, nil
}

func (c *Credentials) NewUser() *User {
	return &User{
		ID:    c.ID,
		Email: c.Email,
	}
}

func (c *Credentials) NewEmail() *UserEmail {
	return &UserEmail{
		ID:     c.Email,
		UserID: c.ID,
		Token:  uuid.New().String(),
	}
}

func (c *Credentials) NewPassword() *UserPassword {
	val, _ := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.MinCost)
	return &UserPassword{
		ID:    c.ID,
		Value: val,
	}
}

type User struct {
	ID           ulid.ULID `json:"id" orm:"pk"`
	Email        string    `json:"email"`
	UserEmail    `json:"-" orm:"fk:Email"`
	UserPassword `json:"-" orm:"fk:ID"`
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
