package contact

import (
	"bytelyon-functions/internal/model"
	"bytelyon-functions/test"
	"context"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestPost(t *testing.T) {

	ctx := context.Background()
	req := test.NewRequest(t).Post(model.Contact{
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
		Value: gofakeit.Sentence(10),
	})

	res, err := Handler(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
}
