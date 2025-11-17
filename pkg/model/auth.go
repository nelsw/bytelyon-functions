package model

import (
	"encoding/base64"
	"errors"
	"strings"
)

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewBasicAuth(s string) (*Auth, error) {

	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(b), ":")
	if len(parts) != 2 {
		return nil, errors.New("invalid basic token; must be base64 encoded '<email>:<password>'")
	}

	return &Auth{
		Username: parts[0],
		Password: parts[1],
	}, nil
}
