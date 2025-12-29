package main

import (
	"bytelyon-functions/pkg/api"
	"bytelyon-functions/pkg/model"
	"net/http"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
)

func init() {
	api.InitLogger()
}

func DemoUser() model.User {
	return model.User{ID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0")}
}

func Test_Handler(t *testing.T) {
	t.Setenv("S3_BUCKET", "bytelyon-db-test")
	user := model.User{ID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0")}

	req := api.NewRequest().
		WithUser(user).
		Get()

	res, _ := Handler(req)

	assert.Equal(t, res.StatusCode, http.StatusOK)
	t.Log(res.Body)
}
