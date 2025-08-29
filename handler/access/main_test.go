package main

import (
	"bytelyon-functions/internal/entity"
	"bytelyon-functions/internal/model/id"
	"bytelyon-functions/internal/model/user"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/brianvoe/gofakeit/v7"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	t.Setenv("APP_MODE", "local")
	u := user.User{
		ID:    id.NewULID(),
		Email: gofakeit.Email(),
	}
	e := user.Email{
		ID:     u.Email,
		UserID: u.ID,
	}
	plainTextPassword := gofakeit.Password(true, true, true, true, true, 8)
	b := []byte(plainTextPassword)
	v, _ := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)
	p := user.Password{
		ID:    u.ID,
		Value: v,
	}
	_ = entity.New().Value(&u).Save()
	_ = entity.New().Value(&e).Save()
	_ = entity.New().Value(&p).Save()

	data := []byte(fmt.Sprintf("%s:%s", u.Email, plainTextPassword))
	ctx := context.Background()
	req := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: http.MethodPost,
			},
		},
		Headers: map[string]string{
			"authorization": "Basic " + base64.StdEncoding.EncodeToString(data),
		},
		QueryStringParameters: map[string]string{
			"login": "true",
		},
	}

	res, _ := handler(ctx, req)
	if res.StatusCode != http.StatusOK {
		t.Fail()
	}

}
