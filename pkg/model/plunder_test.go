package model

import (
	"encoding/json"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestPlunder_Find(t *testing.T) {
	user := MakeDemoUser()
	p := NewPlunder(&user, ulid.MustParse("01KB0P89VRHZBYA68ZGA4R3HMW"))
	err := p.Find()

	assert.NoError(t, err)

	log.Debug().EmbedObject(p).Send()

}

func TestPlunder_FindAll(t *testing.T) {
	user := MakeDemoUser()
	p := NewPlunder(&user)
	arr, err := p.FindAll()
	assert.NoError(t, err)
	assert.NotEmpty(t, arr)
	for _, v := range arr {
		log.Debug().EmbedObject(v).Send()
	}
}

func TestPlunder_Work(t *testing.T) {
	user := MakeDemoUser()
	b, _ := json.Marshal(map[string]any{
		"target": "ev fire blankets",
		"follow": []string{"li-fire.com"},
	})
	p, err := NewPlunder(&user).Create(b)
	assert.NoError(t, err)
	p.Work()
}
