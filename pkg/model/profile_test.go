package model

import (
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestProfile_Save_And_Find(t *testing.T) {

	user := MakeDemoUser()
	expectedName := gofakeit.Name()

	transient := NewProfile(&user)
	transient.Name = expectedName
	b, _ := json.Marshal(&transient)
	persisted, err := transient.Create(b)

	assert.NoError(t, err)
	assert.Equal(t, expectedName, persisted.Name)

	findProxy := NewProfile(&user)
	_, err = findProxy.Find()

	assert.NoError(t, err)
	assert.Equal(t, expectedName, findProxy.Name)
}
