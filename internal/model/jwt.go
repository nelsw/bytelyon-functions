package model

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/lambda"
	"context"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
)

type JWTRequest struct {
	Type  JWTRequestType    `json:"type"`
	Data  map[string]string `json:"data"`
	Token string            `json:"token"`
}

type JWTResponse struct {
	Claims *JWTClaims `json:"claims,omitempty"`
	Token  string     `json:"token,omitempty"`
}

type JWTClaims struct {
	Data map[string]string `json:"data"`
	jwt.RegisteredClaims
}

type JWTRequestType int

const (
	JWTValidation JWTRequestType = iota + 1
	JWTCreation
)

var JWTRequestTypeError = errors.New("invalid JWT request type; must be 1 (validation), 2 (creation)")

func CreateJWT(ctx context.Context, user User) ([]byte, error) {
	return lambda.NewWithContext(ctx).InvokeRequest(ctx, "bytelyon-jwt", app.MustMarshal(JWTRequest{
		Type: JWTCreation,
		Data: map[string]string{"id": user.ID.String()},
	}))
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
	if err != nil {
		return
	}
	var res JWTResponse
	app.MustUnmarshal(out, &res)
	u = User{ID: ulid.MustParse(res.Claims.Data["id"])}
	return
}
