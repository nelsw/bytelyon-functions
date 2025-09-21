package contact

import (
	"bytelyon-functions/test"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestPost(t *testing.T) {
	test.Init(t)
	req := test.NewRequest(t).Post(Contact{
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
		Value: gofakeit.Sentence(10),
	})

	res, _ := Handler(test.CTX, req)

	assert.Equal(t, res.StatusCode, http.StatusOK)
}
