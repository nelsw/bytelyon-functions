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

	res, _ := Handler(api.
		NewRequest().
		User(DemoUser()).
		Method(http.MethodPost).
		Data(model.Job{
			Type: model.NewsJobType,
			Frequency: model.Frequency{
				Unit:  model.Hour,
				Value: 1,
			},
			Name:     gofakeit.Word(),
			Desc:     gofakeit.Sentence(10),
			Keywords: []string{"GMC", "Sierra 1500"},
		}).
		Build())

	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func Test_GET(t *testing.T) {

	res, _ := Handler(api.
		NewRequest().
		User(DemoUser()).
		Method(http.MethodGet).
		Build())

	assert.Equal(t, res.StatusCode, http.StatusOK)
}
