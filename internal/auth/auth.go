package auth

import (
	"bytelyon-functions/internal/db"
	"bytelyon-functions/internal/model"
	"bytelyon-functions/pkg/service/lambda"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

func Login(ctx context.Context, token string) ([]byte, error) {
	if strings.HasPrefix(token, "Basic ") {
		token = strings.TrimPrefix(token, "Basic ")
	}

	b, _ := base64.StdEncoding.DecodeString(token)
	parts := strings.Split(string(b), ":")
	if len(parts) != 2 {
		return nil, errors.New("invalid token")
	}

	email := model.Email{ID: parts[0]}
	if err := db.FindOne(email.Path(), &email); err != nil {
		return nil, err
	}

	pass := model.Password{UserID: email.UserID, Text: parts[1]}
	if err := db.FindOne(pass.Path(), &pass); err != nil {
		return nil, err
	} else if err = pass.Compare(); err != nil {
		return nil, err
	}

	payload, _ := json.Marshal(map[string]any{
		"type": 2,
		"data": map[string]any{"id": email.UserID},
	})

	return lambda.New(ctx).InvokeRequest(ctx, "bytelyon-jwt", payload)
}

func SignUp(ctx context.Context, token string) ([]byte, error) {
	return nil, nil
}
