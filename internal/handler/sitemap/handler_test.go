package main

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/model"
	"bytelyon-functions/test"
	"context"
	"testing"
)

func Test_Handler(t *testing.T) {
	test.Init(t)
	u := model.User{ID: app.NewUlid()}
	Handler(context.Background(), Request{
		UserID: u.ID,
		Depth:  5,
		URL:    "https://www.Li-Fire.com",
	})
}
