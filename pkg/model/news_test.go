package model

import (
	logger "bytelyon-functions/pkg"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestNews_FindAll_News(t *testing.T) {

	user := MakeDemoUser()
	news := NewNews(&user)

	arr, err := news.FindAll()

	assert.NoError(t, err)
	assert.NotEmpty(t, arr)
	t.Setenv("APP_MODE", "test")
	logger.Init()
	for _, n := range arr {
		log.Debug().EmbedObject(n).Msg("news")
	}
}
