package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Data any `json:"data"`
	jwt.RegisteredClaims
}

type Request struct {
	Type  RequestType `json:"type"`
	Data  any         `json:"data"`
	Token string      `json:"token"`
}

type Response struct {
	Claims []byte `json:"claims"`
	Token  string `json:"token"`
	Error  error  `json:"error"`
}

type RequestType int

const (
	Validation RequestType = iota + 1
	Creation
)

type invalidTypeError struct{}

func (invalidTypeError) Error() string {
	return "invalid request type; must be 1 (validation), 2 (creation)"
}

var InvalidRequestType error = invalidTypeError{}

func handler(req Request) Response {

	if req.Type == Validation {

		tkn, err := jwt.ParseWithClaims(req.Token, &Claims{}, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			return Response{Error: err}
		}

		var res Response
		res.Claims, res.Error = json.Marshal(tkn.Claims.(*Claims))
		return res
	}

	if req.Type == Creation {
		tkn, err := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
			Data: req.Data,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    os.Getenv("APP_NAME"),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ID:        uuid.NewString(),
			},
		}).SignedString([]byte(os.Getenv("JWT_SECRET")))

		return Response{Token: tkn, Error: err}
	}

	return Response{Error: InvalidRequestType}
}

func main() {
	lambda.Start(handler)
}
