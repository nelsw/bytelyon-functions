package model

import (
	logger "bytelyon-functions/pkg"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestSearches(t *testing.T) {
	logger.Init()
	t.Setenv("APP_MODE", "test")
	searches, err := NewDemoUser().Searches()
	assert.NoError(t, err)
	assert.NotEmpty(t, searches)

	for _, search := range searches {
		for _, page := range search.Pages {
			log.Debug().EmbedObject(page).Msg("page")
		}
	}
}

func TestSitemaps(t *testing.T) {
	t.Setenv("APP_MODE", "test")
	logger.Init()
	sitemaps, err := NewDemoUser().Sitemaps()
	assert.NoError(t, err)
	assert.NotEmpty(t, sitemaps)
	for _, sitemap := range sitemaps {
		log.Debug().EmbedObject(sitemap).Msg("sitemaps")
	}
}
