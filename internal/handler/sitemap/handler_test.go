package main

import (
	"bytelyon-functions/internal/model"
	"bytelyon-functions/test"
	"context"
	"github.com/oklog/ulid/v2"
	"testing"
)

func Test_Handler(t *testing.T) {
	test.Init(t)
	u := model.User{ID: ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0")}
	Handler(context.Background(), Request{
		UserID: u.ID,
		Depth:  5,
		URL:    "https://www.Li-Fire.com",
	})
}
