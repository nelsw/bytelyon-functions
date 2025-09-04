package model

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type JWTRequest struct {
	Type  JWTRequestType `json:"type"`
	Data  any            `json:"data"`
	Token string         `json:"token"`
}

type JWTResponse struct {
	Claims []byte `json:"claims"`
	Token  string `json:"token"`
}

type JWTClaims struct {
	Data any `json:"data"`
	jwt.RegisteredClaims
}

type JWTRequestType int

const (
	JWTValidation JWTRequestType = iota + 1
	JWTCreation
)

var JWTRequestTypeError = errors.New("invalid JWT request type; must be 1 (validation), 2 (creation)")
