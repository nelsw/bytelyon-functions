package model

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func NewUlid(a ...any) ulid.ULID {

	fn := func(t time.Time) ulid.ULID {
		entropy := rand.New(rand.NewSource(t.UnixNano()))
		ms := ulid.Timestamp(t)
		u, err := ulid.New(ms, entropy)
		if err != nil {
			u = ulid.Make()
		}
		return u
	}

	if len(a) > 0 {
		switch t := a[0].(type) {
		case time.Time:
			return fn(t)
		case string:
			return ulid.MustParse(t)
		}
	}

	return fn(time.Now())
}
