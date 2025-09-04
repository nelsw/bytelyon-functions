package tmp

import (
	"bytelyon-functions/internal/app"

	"github.com/aws/aws-lambda-go/events"
)

func Handler(req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	app.LogURLRequest(req)

	return app.Marshall(req)
}
