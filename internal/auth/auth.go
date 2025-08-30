package auth

import (
	"bytelyon-functions/internal/model"
	"context"
	"encoding/base64"
	"errors"
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

func Validate(ctx context.Context, s string) error {
	if strings.HasPrefix(s, "Bearer ") {
		s = strings.TrimPrefix(s, "Bearer ")
	}

	_, err := jwt.ParseWithClaims(s, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	return err
}

func Login(ctx context.Context, token string) (string, error) {
	if strings.HasPrefix(token, "Basic ") {
		token = strings.TrimPrefix(token, "Basic ")
	}

	b, _ := base64.StdEncoding.DecodeString(token)
	parts := strings.Split(string(b), ":")
	if len(parts) != 2 {
		return "", errors.New("invalid token")
	}

	email := model.Email{ID: parts[0]}
	if err := email.Validate(); err != nil {
		return "", err
	} else if err = model.PasswordText(parts[1]).Validate(); err != nil {
		return "", err
	}

	// todo - get email with metadata and confirm password

	return "", nil
}

func SignUp(ctx context.Context, token string) (string, error) {
	return "", nil
}

func newToken(data interface{}) string {
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
