package auth

import (
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Data interface{} `json:"data"`
	jwt.RegisteredClaims
}

func NewToken(data interface{}) string {
	tkn, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ByteLyon",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}).SignedString([]byte(os.Getenv("JWT_SECRET")))
	return tkn
}

func Validate(s string) (*Claims, error) {

	if strings.HasPrefix(s, "Bearer ") {
		s = strings.TrimPrefix(s, "Bearer ")
	}

	token, err := jwt.ParseWithClaims(s, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims.(*Claims), nil

}
