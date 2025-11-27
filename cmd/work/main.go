package main

import (
	"bytelyon-functions/pkg/model"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
)

func Handler() {
	log.Info().Msg("Doing work...")
	model.DoWork()
	log.Info().Msg("Done.")
}

func main() {
	lambda.Start(Handler)
}
