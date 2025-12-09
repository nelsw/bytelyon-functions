package model

import (
	logger "bytelyon-functions/pkg"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestSearches(t *testing.T) {
	t.Setenv("APP_MODE", "test")
	logger.Init()
	searches, err := NewDemoUser().Searches()
	assert.NoError(t, err)
	assert.NotEmpty(t, searches)

	for _, search := range searches {
		for _, page := range search.Pages {
			log.Debug().EmbedObject(page).Msg("page")
		}
	}
}
