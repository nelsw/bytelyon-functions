package main

import (
	"bytelyon-functions/pkg/service/lambda"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {

	var role, action, name string

	flag.StringVar(&role, "role", "", "The ƒ role to use")
	flag.StringVar(&action, "action", "", "The ƒ action to execute")
	flag.StringVar(&name, "name", "", "The ƒ name to modify")
	flag.Parse()

	function := "bytelyon-" + name

	log.Info().
		Str("action", action).
		Str("handler", name).
		Str("function", function).
		Msg("Handling ƒ request")

	ctx := context.Background()

	if action == "delete" {
		lambda.NewClient(ctx).Delete(ctx, function)
		log.Info().Msgf("Delete function %s success", name)
		return
	}

	if action == "publish" {
		lambda.NewClient(ctx).Publish(ctx, function)
		log.Info().Msgf("Publish function %s success", name)
		return
	}

	if action == "create" || action == "update" {
		env, _ := godotenv.Read(fmt.Sprintf(".handler/%s/.env", name))
		zip := BuildZip(name)
		if action == "create" {
			lambda.NewClient(ctx).Create(ctx, function, role, zip, env)
			log.Info().Msgf("Create function %s success", name)
		} else {
			lambda.NewClient(ctx).Update(ctx, function, zip, env)
			log.Info().Msgf("Update function %s success", name)
		}
		Cleanup()
		return
	}

	log.Warn().Msgf("Unsupported action '%s'", action)
}

func BuildZip(name string) []byte {

	f := "GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap ./handler/%s/main.go"
	if _, err := exec.Command("sh", "-c", fmt.Sprintf(f, name)).Output(); err != nil {
		log.Err(err).Msg("Failed to build zip")
	}
	log.Info().Msgf("Built zip file for ƒ %s", name)

	if _, err := exec.Command("zip", "-9", "-r", "main.zip", "bootstrap").Output(); err != nil {
		log.Err(err).Msg("Failed to build zip")
	}
	log.Info().Msgf("Zipped file for ƒ %s", name)

	b, err := os.ReadFile("main.zip")
	if err != nil {
		log.Err(err).Msg("Failed to build zip")
	}
	log.Info().Msgf("Read file bytes for ƒ %s", name)

	return b
}

func Cleanup() {
	_ = os.Remove("bootstrap")
	_ = os.Remove("main.zip")
	log.Info().Msg("Cleanup finished")
}
