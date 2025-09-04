package model

import (
	"bytelyon-functions/internal/util"
	lamb "bytelyon-functions/pkg/service/lambda"
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
)

type JWTRequest struct {
	Type  JWTRequestType    `json:"type"`
	Data  map[string]string `json:"data"`
	Token string            `json:"token"`
}

type JWTResponse struct {
	Claims *JWTClaims `json:"claims"`
	Token  string     `json:"token"`
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
	return lamb.NewWithContext(ctx).InvokeRequest(ctx, "bytelyon-jwt", util.MustMarshal(JWTRequest{
		Type: JWTCreation,
		Data: map[string]string{"id": user.ID.String()},
	}))
}

func ValidateJWT(ctx context.Context, tkn string) (User, error) {
	out, err := lamb.NewWithContext(ctx).InvokeRequest(ctx, "bytelyon-jwt", util.MustMarshal(JWTRequest{
		Type:  JWTValidation,
		Token: tkn,
	}))
	var u User
	if err == nil {
		var claims JWTClaims
		util.MustUnmarshal(out, &claims)
		id := ulid.MustParse(claims.Data["id"])
		u = User{ID: id}
	}
	return u, err
}
