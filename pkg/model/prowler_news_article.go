package model

import (
	"bytelyon-functions/pkg/util"

	"github.com/oklog/ulid/v2"
)

type NewsItem struct {
	UserID    ulid.ULID `json:"user_id"`
	ProwlerID ulid.ULID `json:"prowler_id"`
	ProwlID   ulid.ULID `json:"prowl_id"`
	ID        ulid.ULID `json:"id"`
	URL       string    `json:"url" xml:"link"`
	Title     string    `json:"title" xml:"title"`
	Source    string    `json:"source"`
	Published *DateTime `json:"-" xml:"pubDate"`
	Sauce     *struct {
		URL   string `json:"-" xml:"url,attr"`
		Value string `json:"-" xml:",chardata"`
	} `json:"-" xml:"source"`
}

func (n NewsItem) String() string {
	return util.Path("user", n.UserID, "prowler", NewsProwlType, n.ProwlerID, "prowl", n.ProwlID, "item", n.ID)
}
