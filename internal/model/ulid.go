package model

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

type ULID struct {
	ID ulid.ULID `json:"id"`
}

func NewULID() ulid.ULID {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	id, err := ulid.New(ms, entropy)
	if err != nil {
		id = ulid.Make()
	}
	return id
}
