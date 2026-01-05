package db

import (
	"errors"
	"strings"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func MagicDelete(userID, entityID ulid.ULID) error {

	log.Debug().Msg("Hold onto your butts ...")

	DB := NewS3()

	keys, err := DB.Keys("user/" + userID.String())
	if err != nil {
		return err
	}

	for _, key := range keys {
		if strings.HasSuffix(key, entityID.String()+".json") ||
			strings.HasSuffix(key, entityID.String()+".html") ||
			strings.HasSuffix(key, entityID.String()+".png") {

			log.Debug().
				Str("key", key).
				Msg("Found a key to get magical with")

			err = errors.Join(DB.Delete(key))
		}
	}

	return err
}
