package sitemap

import (
	"bytelyon-functions/internal/handler/jwt"
	"bytelyon-functions/test"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Handler(t *testing.T) {

	req := test.
		NewRequest(t).
		Bearer(jwt.CreateString(test.CTX, test.FakeUser())).
		Post(Request{URL: "https://ubicquia.com"})

	res, _ := Handler(context.Background(), req)

	assert.Equal(t, res.StatusCode, 200)
}
