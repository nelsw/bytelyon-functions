package model

import "github.com/oklog/ulid/v2"

type Search struct {
	*User     `json:"-"`
	ID        ulid.ULID `json:"id"`
	Query     string
	FollowAds bool     `json:"follow_ads"`
	IgnoreAds []string `json:"except_ads"`
}
