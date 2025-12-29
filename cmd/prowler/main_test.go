package main

import (
	"bytelyon-functions/pkg/api"
	"bytelyon-functions/pkg/model"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Get_Search(t *testing.T) {
	t.Setenv("S3_BUCKET", "bytelyon-db-test")
	req := api.
		NewRequest().
		WithUser(model.User{ID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0")}).
		WithParam("type", model.SearchProwlerType.String()).
		Get()

	res, _ := Handler(req)

	assert.Equal(t, res.StatusCode, http.StatusOK)

	var a []any
	_ = json.Unmarshal([]byte(res.Body), &a)
	b, _ := json.MarshalIndent(a, "", "\t")
	t.Log(string(b))
}

func TestHandler_Get_Sitemap(t *testing.T) {
	t.Setenv("S3_BUCKET", "bytelyon-db-test")
	req := api.
		NewRequest().
		WithUser(model.User{ID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0")}).
		WithParam("type", model.SitemapProwlerType.String()).
		Get()

	res, _ := Handler(req)

	assert.Equal(t, res.StatusCode, http.StatusOK)

	var a []any
	_ = json.Unmarshal([]byte(res.Body), &a)
	b, _ := json.MarshalIndent(a, "", "\t")
	t.Log(string(b))
}

func TestHandler_Get_News(t *testing.T) {
	t.Setenv("S3_BUCKET", "bytelyon-db-test")
	req := api.
		NewRequest().
		WithUser(model.User{ID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0")}).
		WithParam("type", model.NewsProwlerType.String()).
		Get()

	res, _ := Handler(req)

	assert.Equal(t, res.StatusCode, http.StatusOK)
	var a []any
	_ = json.Unmarshal([]byte(res.Body), &a)
	b, _ := json.MarshalIndent(a, "", "\t")
	t.Log(string(b))
}
