package main

import (
	api2 "bytelyon-functions/pkg/api"
	model2 "bytelyon-functions/pkg/model"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
)

func init() {
	api2.InitLogger()
}

func DemoUser() model2.User {
	return model2.User{ID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0")}
}

func Test_POST(t *testing.T) {

	req := api2.NewRequest().
		WithUser(DemoUser()).
		WithData(model2.Job{
			Type: model2.NewsJobType,
			Frequency: model2.Frequency{
				Unit:  model2.Hour,
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

	req := api2.NewRequest().
		WithUser(DemoUser()).
		Get()

	res, _ := Handler(req)

	assert.Equal(t, res.StatusCode, http.StatusOK)
}
