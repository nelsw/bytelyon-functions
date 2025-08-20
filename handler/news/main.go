package news

import (
	"bytelyon-functions/pkg/api"
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	// get all jobs
	// get all work
	// get all articles

	// force job
	// create job
	// update job

	// delete job

	// delete articles
	// batch delete articles
	return api.Response(http.StatusOK, "")
}

func main() {
	lambda.Start(handler)
}
