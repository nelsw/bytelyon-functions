package model

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestProfile_Save(t *testing.T) {
	user := MakeDemoUser()
	oldProfile := NewProfile(&user)
	newProfile := Profile{Name: gofakeit.Name()}

	assert.NoError(t, oldProfile.Hydrate(newProfile))
	assert.NoError(t, oldProfile.Validate())
	assert.NoError(t, oldProfile.Save())
	assert.Equal(t, oldProfile.Name, newProfile.Name)
}

func TestProfile_Find(t *testing.T) {
	user := MakeDemoUser()
	oldProfile := NewProfile(&user)
	newProfile := Profile{Name: gofakeit.Name()}

	assert.NoError(t, oldProfile.Hydrate(newProfile))
	assert.NoError(t, oldProfile.Validate())
	assert.NoError(t, oldProfile.Save())
	assert.NoError(t, oldProfile.Find())
	assert.Equal(t, oldProfile.Name, newProfile.Name)
}
