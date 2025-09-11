package model

import (
	"bytelyon-functions/internal/client/s3"
	"bytelyon-functions/test"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindAllUsers(t *testing.T) {
	test.Init(t)
	db := s3.New(test.CTX)
	users, err := FindAllUsers(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, users)
	b, _ := json.MarshalIndent(users, "", "  ")
	t.Log(string(b))
}
