package util

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/oklog/ulid/v2"
)

func Ptr[T any](a T) *T { return &a }

func Between[T int | float64](min, max T) T {
	return T(rand.Intn(int(max)-int(min)) + int(min))
}

func Path(a ...any) string { return Join("/", a...) }

func Join(sep string, a ...any) string {
	var strs []string
	for _, v := range a {
		switch t := v.(type) {
		case string:
			strs = append(strs, t)
		case ulid.ULID:
			if !t.IsZero() {
				strs = append(strs, t.String())
			}
		case fmt.Stringer:
			strs = append(strs, t.String())
		}
	}
	return strings.Join(strs, sep)
}
