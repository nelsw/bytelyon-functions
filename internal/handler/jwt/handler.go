package jwt

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/lambda"
	"bytelyon-functions/internal/model"
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Request struct {
	Type  RequestType `json:"type"`
	Data  model.User  `json:"data"`
	Token string      `json:"token"`
}

type Response struct {
	Claims *Claims `json:"claims,omitempty"`
	Token  string  `json:"token,omitempty"`
}

type Claims struct {
	Data model.User `json:"data"`
	jwt.RegisteredClaims
}

type RequestType int

const (
	Validation RequestType = iota + 1
	Creation
)

const name = "bytelyon-jwt"

var (
	typeError = errors.New("invalid JWT request type; must be 1 (validation), 2 (creation)")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	jwtIssuer = os.Getenv("APP_NAME")
)

func Handler(req Request) (res Response, err error) {

	log.Info().Any("request", req).Send()

	if req.Type == Validation {
		var tkn *jwt.Token
		if tkn, err = jwt.ParseWithClaims(req.Token, &Claims{}, func(token *jwt.Token) (any, error) {
			return jwtSecret, nil
		}); err == nil {
			res.Claims = tkn.Claims.(*Claims)
		}
		log.Err(err).Any("response", res).Msg("validate JWT")
		return
	}

	if req.Type == Creation {
		res.Token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
			Data: req.Data,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    jwtIssuer,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ID:        uuid.NewString(),
			},
		}).SignedString(jwtSecret)
		log.Err(err).Any("response", res).Msg("create JWT")
		return
	}

	err = typeError
	log.Err(err).Send()
	return
}

func Validate(ctx context.Context, tkn string) (u model.User, err error) {
	var out []byte
	out, err = lambda.New(ctx).Request(ctx, name, Request{
		Type:  Validation,
		Token: strings.TrimPrefix(tkn, "Bearer "),
	})
	if strings.Contains(string(out), "error") {
		err = errors.Join(err, errors.New(string(out)))
	}
	if err == nil {
		var res Response
		app.MustUnmarshal(out, &res)
		u = res.Claims.Data
	}
	return
}

func Create(ctx context.Context, user model.User) (out []byte, err error) {
	out, err = lambda.New(ctx).Request(ctx, name, Request{
		Type: Creation,
		Data: user,
	})
	if strings.Contains(string(out), "error") {
		err = errors.Join(err, errors.New(string(out)))
	}
	return
}

func CreateString(ctx context.Context, user model.User) string {
	out, err := Create(ctx, user)
	if err != nil {
		log.Panic().Err(err).Send()
	}
	var res Response
	app.MustUnmarshal(out, &res)
	return res.Token
}
