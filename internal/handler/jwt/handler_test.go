package jwt

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

// expired token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7ImlkIjoiMDFLNDhQQzBCSzEzQldWMkNHV0ZQOFFRSDAifSwiaXNzIjoiQnl0ZUx5b24iLCJleHAiOjE3NTcwMTgxNDMsIm5iZiI6MTc1NzAxNjM0MywiaWF0IjoxNzU3MDE2MzQzLCJqdGkiOiIwYWVlNDdjMy03YTQ2LTRjYmQtYTdhYy1jNzQ2NjBmODg0MjQifQ.04abFJOZf-qB1C-C2y7Pjj4c2krkAyxCDZy7SK7p3Y4

func TestBadType(t *testing.T) {
	t.Setenv("JWT_SECRET", "a-string-secret-at-least-256-bits-long")
	res, err := Handler(model.JWTRequest{})
	assert.Equal(t, res, model.JWTResponse{})
	assert.ErrorIs(t, err, model.JWTRequestTypeError)
}

func TestOK(t *testing.T) {
	t.Setenv("JWT_SECRET", "a-string-secret-at-least-256-bits-long")

	data := model.User{ID: app.NewUlid()}
	res, err := Handler(model.JWTRequest{
		Type: model.JWTCreation,
		Data: data,
	})
	assert.NoError(t, err)

	res, err = Handler(model.JWTRequest{
		Type:  model.JWTValidation,
		Token: res.Token,
	})
	assert.NoError(t, err)

	assert.NoError(t, err)
	assert.Equal(t, res.Claims.Data, data)
}
