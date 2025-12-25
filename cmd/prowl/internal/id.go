package internal

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func NewULID(s ...any) ulid.ULID {

	f := func(t time.Time) ulid.ULID {
		utc := t.UTC()
		ms := ulid.Timestamp(utc)
		entropy := rand.New(rand.NewSource(utc.UnixNano()))
		return ulid.MustNew(ms, entropy)
	}

	if len(s) == 0 {
		return f(time.Now().UTC())
	}

	if id, ok := s[0].(string); ok {
		return ulid.MustParse(id)
	}

	return f(s[0].(time.Time))
}
