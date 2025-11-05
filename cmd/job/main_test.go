package main

import (
	"bytelyon-functions/internal/api"
	"bytelyon-functions/internal/model"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
)

func init() {
	api.InitLogger()
}

func DemoUser() model.User {
	return model.User{ID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0")}
}

func Test_POST(t *testing.T) {

	req := api.
		NewRequest().
		WithUser(DemoUser()).
		WithData(model.Job{
			Type: model.NewsJobType,
			Frequency: model.Frequency{
				Unit:  model.Hour,
				Value: 1,
			},
			Name:     gofakeit.Word(),
			Desc:     gofakeit.Sentence(10),
			Keywords: []string{"GMC", "Sierra 1500"},
		}).Post()

	res, _ := Handler(req)

	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func Test_GET(t *testing.T) {

	req := api.
		NewRequest().
		WithUser(DemoUser()).
		Get()

	res, _ := Handler(req)

	assert.Equal(t, res.StatusCode, http.StatusOK)
}
