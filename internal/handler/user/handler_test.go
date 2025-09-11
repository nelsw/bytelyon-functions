package user

import (
	"bytelyon-functions/internal/client/s3"
	"bytelyon-functions/internal/handler/jwt"
	"bytelyon-functions/internal/model"
	"bytelyon-functions/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FindAll_Users(t *testing.T) {
	test.Init(t)

	db := s3.New(test.CTX)

	users, err := model.FindAllUsers(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, users)

	user := users[0]

	req := test.
		NewRequest(t).
		Bearer(jwt.CreateString(test.CTX, user)).
		Query("delimiter", model.UserPath).
		Get()

	res, _ := Handler(test.CTX, req)

	assert.Equal(t, res.StatusCode, 200)
}

func Test_FindAll_Sitemaps(t *testing.T) {
	test.Init(t)

	db := s3.New(test.CTX)

	users, err := model.FindAllUsers(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, users)

	user := users[0]

	req := test.
		NewRequest(t).
		Bearer(jwt.CreateString(test.CTX, user)).
		Query("delimiter", "sitemap").
		Get()

	res, _ := Handler(test.CTX, req)

	assert.Equal(t, res.StatusCode, 200)
}

func Test_FindAll_Jobs(t *testing.T) {
	test.Init(t)

	db := s3.New(test.CTX)

	users, err := model.FindAllUsers(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, users)

	user := users[0]

	req := test.
		NewRequest(t).
		Bearer(jwt.CreateString(test.CTX, user)).
		Query("delimiter", "job").
		Get()

	res, _ := Handler(test.CTX, req)

	assert.Equal(t, res.StatusCode, 200)
}

func Test_FindAll_Items(t *testing.T) {
	test.Init(t)

	db := s3.New(test.CTX)

	users, err := model.FindAllUsers(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, users)

	user := users[0]

	req := test.
		NewRequest(t).
		Bearer(jwt.CreateString(test.CTX, user)).
		Query("delimiter", "item").
		Get()

	res, _ := Handler(test.CTX, req)

	assert.Equal(t, res.StatusCode, 200)
}
