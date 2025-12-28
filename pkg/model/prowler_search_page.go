package model

import (
	"bytelyon-functions/pkg/db"
	"bytelyon-functions/pkg/util"
	"errors"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

// SearchPage todo - refactor to Page
type SearchPage struct {
	UserID     ulid.ULID `json:"user_id"`
	ProwlerID  ulid.ULID `json:"prowler_id"`
	ProwlID    ulid.ULID `json:"prowl_id"`
	ID         ulid.ULID `json:"id"`
	Title      string    `json:"title"`
	URL        string    `json:"url"`
	Domain     string    `json:"domain"`
	Data       any       `json:"data"`
	Content    string    `json:"-"`
	Screenshot []byte    `json:"-"`
}

func (p *SearchPage) String() string {
	return util.Path("user", p.UserID, "prowler", SearchProwlerType, p.ProwlerID, "prowl", p.ProwlID, "page", p.ID)
}

func (p *SearchPage) Save() {

	DB := db.NewS3()

	var err error
	err = errors.Join(err, DB.Put(p.String()+"/screenshot.png", p.Screenshot))
	err = errors.Join(err, DB.Put(p.String()+"/content.html", []byte(p.Content)))
	err = errors.Join(err, db.Save(p))

	if err != nil {
		log.Warn().Err(err).Stringer("key", p).Msg("Prowler - failed to save page")
	} else {
		log.Info().Stringer("key", p).Msg("Prowler - Saved page")
	}
}
