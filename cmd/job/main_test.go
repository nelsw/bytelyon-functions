package main

import (
	"bytelyon-functions/internal/api"
	"bytelyon-functions/internal/model"
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

	req := api.
		NewRequest().
		User(DemoUser()).
		Method(http.MethodPost).
		Data(model.Job{
			Type: model.NewsJobType,
			Frequency: model.Frequency{
				Unit:  model.Hour,
				Value: 1,
			},
			Name:     "Test Name",
			Desc:     "Test Desc",
			Keywords: []string{"GMC", "Sierra 1500"},
		}).
		Build()

	res, _ := Handler(req)

	assert.Equal(t, res.StatusCode, http.StatusOK)
}
