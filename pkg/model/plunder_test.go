package model

import (
	logger "bytelyon-functions/pkg"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestPlunder_Find(t *testing.T) {
	t.Setenv("APP_MODE", "test")
	logger.Init()
	user := MakeDemoUser()
	p := NewPlunder(&user, ulid.MustParse("01KBK75WKBNWQJS0R11G6XV8YG"))
	err := p.Find()

	assert.NoError(t, err)

	log.Debug().EmbedObject(p).Msg("plunder")
	for _, v := range p.Loot {
		log.Debug().EmbedObject(v).Msg("loot")
	}
}

func TestPlunder_FindAll(t *testing.T) {
	t.Setenv("APP_MODE", "test")
	logger.Init()
	user := MakeDemoUser()
	p := NewPlunder(&user)
	arr, err := p.FindAll()
	assert.NoError(t, err)
	assert.NotEmpty(t, arr)
	for _, p = range arr {
		log.Debug().EmbedObject(p).Msg("plunder")
		for _, v := range p.Loot {
			log.Debug().EmbedObject(v).Msg("loot")
		}
	}
}

func TestPlunder_Work(t *testing.T) {
	t.Setenv("APP_MODE", "prod")
	logger.Init()
	user := MakeDemoUser()
	p := NewPlunder(&user, ulid.MustParse("01KBK75WKBNWQJS0R11G6XV8YG"))
	p.Work()
}
