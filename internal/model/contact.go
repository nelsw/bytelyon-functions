package model

import (
	"bytelyon-functions/internal/util"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/oklog/ulid/v2"
)

type Contact struct {
	ID    ulid.ULID `json:"id" fake:"skip"`
	Name  string    `json:"name" fake:"{name}"`
	Email string    `json:"email" fake:"{email}"`
	Value string    `json:"message" fake:"{sentence}"`
}

func NewContact(s string) (*Contact, error) {
	var c Contact
	if err := json.Unmarshal([]byte(s), &c); err != nil {
		return nil, err
	} else if err = c.Validate(); err != nil {
		return nil, err
	}
	c.ID = NewUlid()
	return &c, nil
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

func (c Contact) Key() string {
	return fmt.Sprintf("%s/db/message/contact/unread/%s.json", util.AppMode(), c.ID)
}
