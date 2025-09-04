package repo

import (
	"bytelyon-functions/test"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_GetOne(t *testing.T) {
	ctx := context.Background()
	req := test.NewRequest(t).
		Path("user/01K3Z13PH6JMGJCYF0Z6V166MQ").
		Get()

	res, err := Handler(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 401)
}
