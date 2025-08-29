package credentials

import (
	"bytelyon-functions/internal/model/user"
	"encoding/base64"
	"errors"
	"net/mail"
	"strings"
	"unicode"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

func (c *Credentials) NewUser() *user.User {
	return &user.User{
		ID:    ulid.Make(),
		Email: c.Email,
		Roles: []user.RoleType{user.Basic},
	}
}

func (c *Credentials) NewUserProfile(userID ulid.ULID) *user.Profile {
	return &user.Profile{ID: userID}
}

func (c *Credentials) NewEmail(userID ulid.ULID) *user.Email {
	return &user.Email{
		ID:     c.Email,
		UserID: userID,
		Token:  ulid.Make().String(),
	}
}

func (c *Credentials) NewPassword(userID ulid.ULID) *user.Password {
	val, _ := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.MinCost)
	return &user.Password{
		ID:    userID,
		Value: val,
	}
}
