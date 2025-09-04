package contact

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/s3"
	"bytelyon-functions/internal/model"
	"context"

	"github.com/aws/aws-lambda-go/events"
)

func Handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	app.LogURLRequest(req)

	if app.IsOptions(req) {
		return app.OK()
	}

	if app.IsPost(req) {
		c, err := model.MakeContact(req.Body)
		if err == nil {
			err = s3.NewWithContext(ctx).Put(c.Path(), app.MustMarshal(c))
		}
		if err != nil {
			return app.BadRequest(err)
		}
		return app.Marshall(c)
	}

	return app.NotImplemented(req)
}
