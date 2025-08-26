package bot

import (
	"bytelyon-functions/internal/entity"
	"bytelyon-functions/internal/model/id"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type Work struct {
	ID       ulid.ULID `json:"id"`
	JobID    uuid.UUID `json:"job_id"`
	Root     string    `json:"root"`
	Keywords []string  `json:"keywords"`
	Items    []Item    `json:"items"`
}

func createWorkItems(ID uuid.UUID, t JobType, u string, keys []string) (err error) {

	u = fmt.Sprintf(u, strings.Join(keys, ","))

	var res *http.Response
	if res, err = http.Get(u); err != nil {
		log.Error().Err(err).Str("URL", u).Str("ID", ID.String()).Msg("failed to http.Get url")
		return
	}
	defer res.Body.Close()

	var b []byte
	if b, err = io.ReadAll(res.Body); err != nil {
		log.Error().Err(err).Str("ID", ID.String()).Msg("failed to io.ReadAll response body")
		return
	}

	b = t.Sanitize(b)

	var rss RSS
	if err = xml.Unmarshal(b, &rss); err != nil {
		log.Error().Err(err).Str("ID", ID.String()).Msg("failed to unmarshal xml")
		return
	}

	log.Info().Int("size", len(rss.Channel.Items)).Msg("work items found")

	if len(rss.Channel.Items) == 0 {
		return
	}

	work := Work{
		ID:       id.NewULID(),
		JobID:    ID,
		Items:    rss.Channel.Items,
		Root:     u,
		Keywords: keys,
	}
	if err = entity.New().Value(&work).Save(); err != nil {
		return
	}

	log.Info().Str("workID", work.ID.String()).Msg("work items created")

	return
}
