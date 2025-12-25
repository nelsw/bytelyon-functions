package internal

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func Path(a ...any) string {
	path, size := "", len(a)
	for i, s := range a {
		if id, ok := s.(ulid.ULID); ok && id.IsZero() {
			break
		}
		if path += fmt.Sprint(s); i < size-1 {
			path += "/"
		}
	}
	return path
}

func Ptr[T any](a T) *T { return &a }

func Bytes(a any) []byte {
	switch v := a.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	default:
		b, err := json.Marshal(a)
		if err != nil {
			log.Panic().Err(err).Any("arg", a).Msg("Bytes failed to marshal")
		}
		return b
	}
}

func Between[T int | float64](min, max T) T {
	return T(rand.Intn(int(max)-int(min)) + int(min))
}
