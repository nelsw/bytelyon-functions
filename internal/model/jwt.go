package model

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/lambda"
	"context"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type JWTRequest struct {
	Type  JWTRequestType `json:"type"`
	Data  User           `json:"data"`
	Token string         `json:"token"`
}

type JWTResponse struct {
	Claims *JWTClaims `json:"claims,omitempty"`
	Token  string     `json:"token,omitempty"`
}

type JWTClaims struct {
	Data User `json:"data"`
	jwt.RegisteredClaims
}

type JWTRequestType int

const (
	JWTValidation JWTRequestType = iota + 1
	JWTCreation
)

var JWTRequestTypeError = errors.New("invalid JWT request type; must be 1 (validation), 2 (creation)")

func CreateJWT(ctx context.Context, user User) (out []byte, err error) {
	out, err = lambda.NewWithContext(ctx).InvokeRequest(ctx, "bytelyon-jwt", app.MustMarshal(JWTRequest{
		Type: JWTCreation,
		Data: user,
	}))
	if strings.Contains(string(out), "error") {
		err = errors.Join(err, errors.New(string(out)))
	}
	return
}

func CreateJWTString(ctx context.Context, user User) string {
	out, err := CreateJWT(ctx, user)
	if err != nil {
		log.Panic().Err(err).Send()
	}
	var res JWTResponse
	app.MustUnmarshal(out, &res)
	return res.Token
}

func ValidateJWT(ctx context.Context, tkn string) (u User, err error) {
	var out []byte
	out, err = lambda.NewWithContext(ctx).InvokeRequest(ctx, "bytelyon-jwt", app.MustMarshal(JWTRequest{
		Type:  JWTValidation,
		Token: tkn,
	}))
	if strings.Contains(string(out), "error") {
		err = errors.Join(err, errors.New(string(out)))
	}
	if err == nil {
		var res JWTResponse
		app.MustUnmarshal(out, &res)
		u = res.Claims.Data
	}
	return
}
