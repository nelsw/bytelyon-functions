package jwt

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/model"
	"bytelyon-functions/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

// expired token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7ImlkIjoiMDFLNDhQQzBCSzEzQldWMkNHV0ZQOFFRSDAifSwiaXNzIjoiQnl0ZUx5b24iLCJleHAiOjE3NTcwMTgxNDMsIm5iZiI6MTc1NzAxNjM0MywiaWF0IjoxNzU3MDE2MzQzLCJqdGkiOiIwYWVlNDdjMy03YTQ2LTRjYmQtYTdhYy1jNzQ2NjBmODg0MjQifQ.04abFJOZf-qB1C-C2y7Pjj4c2krkAyxCDZy7SK7p3Y4

func TestBadType(t *testing.T) {
	test.Init(t)
	res, err := Handler(Request{})
	assert.Equal(t, res, Response{})
	assert.ErrorIs(t, err, typeError)
}

func TestOK(t *testing.T) {
	test.Init(t)

	data := model.User{ID: app.NewUlid()}
	res, err := Handler(Request{
		Type: Creation,
		Data: data,
	})
	assert.NoError(t, err)

	res, err = Handler(Request{
		Type:  Validation,
		Token: res.Token,
	})
	assert.NoError(t, err)

	assert.NoError(t, err)
	assert.Equal(t, res.Claims.Data, data)
}
