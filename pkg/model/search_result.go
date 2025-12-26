package model

import "github.com/oklog/ulid/v2"

type SearchResult struct {
	UserID   ulid.ULID `json:"user_id"`
	SearchID ulid.ULID `json:"search_id"`
	ID       ulid.ULID `json:"id"`
	Pages    []*Page   `json:"pages"`
}
