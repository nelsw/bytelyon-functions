package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNews_FindAll_News(t *testing.T) {

	user := MakeDemoUser()
	news := NewNews(&user)

	arr, err := news.FindAll()

	assert.NoError(t, err)
	assert.NotEmpty(t, arr)
}
