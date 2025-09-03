package util

import (
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func Test_IsJSON_True(t *testing.T) {
	b, err := json.Marshal(&map[string]any{
		gofakeit.Word(): gofakeit.Word(),
		gofakeit.Word(): gofakeit.Int64(),
		gofakeit.Word(): gofakeit.Float64(),
	})
	assert.NoError(t, err)

	ok := IsJSON(string(b))
	assert.True(t, ok)
}

func Test_IsJSON_False(t *testing.T) {
	ok := IsJSON("")
	assert.False(t, ok)
}

func Test_First_NotNil(t *testing.T) {
	var exp any
	act := First(exp)
	assert.Equal(t, exp, act)
}

func Test_First_Nil(t *testing.T) {
	act := First(nil)
	assert.Nil(t, act)
}
