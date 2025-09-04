package app

import (
	"errors"
	"os"

	"github.com/rs/zerolog/log"
)

func Mode() string {
	mode := os.Getenv("APP_MODE")
	if mode == "" {
		log.Panic().Err(errors.New("APP_MODE is required")).Send()
	}
	return mode
}

func Bucket() *string {
	s := "bytelyon-db-" + Mode()
	return &s
}
