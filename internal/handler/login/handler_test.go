package login

import (
	"bytelyon-functions/test"
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	t.Setenv("APP_MODE", "test")
	e := "demo@demo.com"
	p := "Demo123!"
	s := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", e, p)))
	ctx := context.Background()
	req := test.NewRequest(t).Path("login").Header("authorization", s).Post()

	res, err := Handler(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 200)
}
