package jwt

import (
	"bytelyon-functions/internal/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func Handler(req model.JWTRequest) (res model.JWTResponse, err error) {

	log.Info().Any("request", req).Send()

	if req.Type == model.JWTValidation {
		var tkn *jwt.Token
		if tkn, err = jwt.ParseWithClaims(req.Token, &model.JWTClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		}); err == nil {
			res.Claims = tkn.Claims.(*model.JWTClaims)
		}
		log.Err(err).Any("response", res).Msg("validate JWT")
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
		log.Err(err).Any("response", res).Msg("create JWT")
		return
	}

	err = model.JWTRequestTypeError
	log.Err(err).Send()
	return
}
