package main

import (
	"bytelyon-functions/internal/model"
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestBadType(t *testing.T) {
	t.Setenv("JWT_SECRET", "a-string-secret-at-least-256-bits-long")
	res, err := handler(model.JWTRequest{})
	assert.Equal(t, res, model.JWTResponse{})
	assert.ErrorIs(t, err, model.JWTRequestTypeError)
}

func TestOK(t *testing.T) {
	t.Setenv("JWT_SECRET", "a-string-secret-at-least-256-bits-long")

	data := map[string]any{"id": gofakeit.UUID()}
	res, err := handler(model.JWTRequest{
		Type: model.JWTCreation,
		Data: data,
	})
	assert.NoError(t, err)

	res, err = handler(model.JWTRequest{
		Type:  model.JWTValidation,
		Token: res.Token,
	})
	assert.NoError(t, err)

	var m map[string]any
	err = json.Unmarshal(res.Claims, &m)
	assert.NoError(t, err)
	assert.Equal(t, m["data"], data)
}
