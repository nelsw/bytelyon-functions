package model

import (
	"bytelyon-functions/internal/client/s3"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type User struct {
	ID ulid.ULID `json:"id"`
}

func NewUser(req events.APIGatewayV2HTTPRequest) (*User, error) {

	log.Info().Any("request", req).Send()

	b, err := json.Marshal(req.RequestContext.Authorizer.Lambda["user"])
	if err != nil {
		return nil, err
	}

	var u User
	if err = json.Unmarshal(b, &u); err != nil {
		return nil, err
	}

	return &u, nil
}

func (u *User) Path() string {
	return "user/" + u.ID.String()
}

func (u *User) Key() string {
	return u.Path() + "/_.json"
}

func FindUser(s string) (*User, error) {

	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(b), ":")
	if len(parts) != 2 {
		return nil, errors.New("invalid basic token; must be base64 encoded '<email>:<password>'")
	}

	var e *Email
	var p *Password
	if e, err = NewEmail(parts[0]); err != nil {
		return nil, err
	} else if p, err = NewPassword(parts[1]); err != nil {
		return nil, err
	}

	db := s3.New()
	if err = e.Find(db); err != nil {
		return nil, err
	}

	if err = p.Find(db, e.User()); err != nil {
		return nil, err
	}

	return e.User(), nil
}
