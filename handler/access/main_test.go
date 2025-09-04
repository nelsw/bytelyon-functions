package main

import (
	"bytelyon-functions/pkg/api"
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {

	e := "demo@demo.com"
	p := "Demo123!"
	s := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", e, p)))
	ctx := context.Background()
	req := api.NewRequest().Path("login").Header("authorization", s).Post()

	res, err := handler(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 200)
}
