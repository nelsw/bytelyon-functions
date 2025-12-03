package main

import (
	"bytelyon-functions/pkg/api"
	"bytelyon-functions/pkg/model"
	"encoding/json"
	"testing"

	"github.com/rs/zerolog/log"
)

func Test(t *testing.T) {
	api.InitLogger()
	req := api.NewRequest().WithUser(model.MakeDemoUser()).Get()
	t.Setenv("APP_MODE", "prod")

	res, _ := Handler(req)

	var arr []model.Plunder
	json.Unmarshal([]byte(res.Body), &arr)

	for _, v := range arr {
		log.Debug().EmbedObject(&v).Send()
		for _, l := range v.Loots {
			log.Debug().EmbedObject(l).Send()
		}
	}
}
