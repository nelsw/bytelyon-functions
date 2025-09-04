package main

import (
	"bytelyon-functions/internal/model"
	"bytelyon-functions/internal/util"
	"bytelyon-functions/pkg/api"
	"bytelyon-functions/pkg/service/s3"
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const bucket = "bytelyon"

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	api.LogURLRequest(req)

	switch {
	case api.IsOptions(req):
		return api.OK()
	case api.IsPost(req):
		c, err := model.NewContact(req.Body)
		if err == nil {
			err = s3.NewWithContext(ctx).Put(ctx, "bytelyon", c.Key(), util.MustMarshal(*c))
		}
		return api.Err(err)
	default:
		return api.NotImplemented(req)
	}
}

func main() {
	lambda.Start(handler)
}
