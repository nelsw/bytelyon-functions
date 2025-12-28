package model

import (
	"bytelyon-functions/internal"
	"encoding/json"
	"errors"
	"fmt"

	//"github.com/nelsw/bytelyon-db"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type Profile struct {
	UserID ulid.ULID `json:"user_id"`
	Name   string    `json:"name"`
}

func NewProfile(u *User) *Profile {
	return &Profile{UserID: u.ID}
}

func (p *Profile) String() string {
	return fmt.Sprintf("user/%s/profile", p.UserID)
}

func (p *Profile) Create(b []byte) (*Profile, error) {
	var v = new(Profile)
	if err := json.Unmarshal(b, v); err != nil {
		log.Err(err).Msg("failed to unmarshal profile")
		return nil, err
	}
	if v.Name == "" {
		return nil, errors.New("name cannot be empty")
	}
	v.UserID = p.UserID
	return v, db.Save(v)
}

func (p *Profile) Find() (*Profile, error) {
	return p, db.Find(p)
}
