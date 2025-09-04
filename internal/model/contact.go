package model

import (
	"bytelyon-functions/internal/app"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/oklog/ulid/v2"
)

type Contact struct {
	ID    ulid.ULID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Value string    `json:"message"`
}

func MakeContact(s string) (c Contact, err error) {
	if err = json.Unmarshal([]byte(s), &c); err != nil {
		return
	} else if err = c.Validate(); err != nil {
		return
	}
	c.ID = app.NewUlid()
	return c, nil
}

func (c Contact) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	} else if c.Email == "" {
		return errors.New("email is required")
	} else if c.Value == "" {
		return errors.New("message is required")
	}
	return nil
}

func (c Contact) Path() string {
	return fmt.Sprintf("message/contact/unread/%s", c.ID)
}
