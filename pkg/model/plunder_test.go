package model

import (
	logger "bytelyon-functions/pkg"
	"encoding/json"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestPlunder_Find(t *testing.T) {
	t.Setenv("APP_MODE", "test")
	logger.Init()
	user := MakeDemoUser()
	p := NewPlunder(&user, ulid.MustParse("01KBH5HA4358EG5W61N4S8RPN6"))
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
	b, _ := json.Marshal(map[string]any{
		"target": "ev fire blankets",
		"follow": []string{"li-fire.com"},
	})
	p, err := NewPlunder(&user).Create(b)
	assert.NoError(t, err)
	p.Work()
}
