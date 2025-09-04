package jwt

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
