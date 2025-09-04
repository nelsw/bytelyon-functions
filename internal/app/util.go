package app

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func MustMarshal(a any) []byte {
	b, err := json.Marshal(&a)
	if err != nil {
		log.Panic().
			Err(err).
			Any("any", a).
			Msg("MustMarshal")
	}
	return b
}

func MustUnmarshal(b []byte, a any) {
	if err := json.Unmarshal(b, &a); err != nil {
		log.Panic().
			Err(err).
			Bytes("bytes", b).
			Any("any", a).
			Msg("MustUnmarshal")
	}
}

func IsJSON(s string) bool {
	if s == "" {
		return false
	}
	var raw json.RawMessage
	return json.Unmarshal([]byte(s), &raw) == nil
}

func NewUlid() ulid.ULID {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	u, err := ulid.New(ms, entropy)
	if err != nil {
		u = ulid.Make()
	}
	return u
}
