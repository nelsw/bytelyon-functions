package sitemap

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/s3"
	"bytelyon-functions/internal/handler/jwt"
	"bytelyon-functions/internal/model"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"

	"github.com/oklog/ulid/v2"
)

func Handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	if app.LogURLRequest(req); app.IsOptions(req) {
		return app.OK()
	}

	var userID ulid.ULID
	if user, err := jwt.Validate(ctx, req.Headers["authorization"]); err != nil {
		return app.Unauthorized(err)
	} else {
		userID = user.ID
	}

	if app.IsPut(req) || app.IsPost(req) {

		var s model.Sitemap
		if err := json.Unmarshal([]byte(req.Body), &s); err != nil {
			return app.BadRequest(err)
		}

		return app.Response(s.Create(s3.New(ctx), userID))
	}

	return app.NotImplemented(req)
}
