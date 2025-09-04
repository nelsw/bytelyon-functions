package main

import (
	"bytelyon-functions/internal/model"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func handler(req model.JWTRequest) (res model.JWTResponse, err error) {

	if req.Type == model.JWTValidation {
		var tkn *jwt.Token
		if tkn, err = jwt.ParseWithClaims(req.Token, &model.JWTClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		}); err == nil {
			res.Claims = tkn.Claims.(*model.JWTClaims)
		}
		return
	}

	if req.Type == model.JWTCreation {
		res.Token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, model.JWTClaims{
			Data: req.Data,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    os.Getenv("APP_NAME"),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ID:        uuid.NewString(),
			},
		}).SignedString([]byte(os.Getenv("JWT_SECRET")))
		return
	}

	err = model.JWTRequestTypeError

	return
}

func main() {
	lambda.Start(handler)
}
