package model

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func NewUlid() ulid.ULID {
	return newUlid(time.Now())
}

func NewUlidFromTime(t time.Time) ulid.ULID {
	return newUlid(t)
}

func newUlid(t time.Time) ulid.ULID {
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	ms := ulid.Timestamp(time.Now())
	u, err := ulid.New(ms, entropy)
	if err != nil {
		u = ulid.Make()
	}
	return u
}
