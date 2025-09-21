package main

import (
	. "bytelyon-functions/internal/handler/tmp"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(Handler)
}
