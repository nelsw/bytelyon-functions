package main

import (
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestBadType(t *testing.T) {
	res := handler(Request{})
	assert.ErrorIs(t, res.Error, InvalidRequestType)
}

func TestOK(t *testing.T) {

	data := map[string]any{"id": gofakeit.UUID()}
	creationResponse := handler(Request{
		Type: Creation,
		Data: data,
	})
	assert.Nil(t, creationResponse.Error)

	validationResponse := handler(Request{
		Type:  Validation,
		Token: creationResponse.Token,
	})
	assert.Nil(t, validationResponse.Error)

	var m map[string]any
	err := json.Unmarshal(validationResponse.Claims, &m)
	assert.Nil(t, err)
	assert.Equal(t, m["data"], data)
}

func TestValidation_Err(t *testing.T) {
	res := handler(Request{
		Type:  Validation,
		Token: "a-bad-token",
	})
	assert.Error(t, res.Error)
}
