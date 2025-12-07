package model

import (
	"github.com/rs/zerolog"
)

type Loots []*Loot

func (ll Loots) MarshalZerologArray(a *zerolog.Array) {
	for _, l := range ll {
		a.Object(l)
	}
}
