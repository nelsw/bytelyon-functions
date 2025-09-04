package job

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/s3"
	"bytelyon-functions/internal/model"
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func Handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	app.LogURLRequest(req)

	if app.IsOptions(req) {
		return app.OK()
	}

	if app.IsPut(req) || app.IsPost(req) {
		user, err := model.ValidateJWT(ctx, req.Headers["authorization"])
		if err != nil {
			return app.Unauthorized(err)
		} else if !strings.Contains("user/"+user.ID.String(), req.RawPath) {
			return app.Forbidden()
		}

		var j model.Job
		if err = json.Unmarshal([]byte(req.Body), &j); err == nil {
			err = j.Validate()
		}

		if err != nil {
			return app.BadRequest(err)
		}

		if j.ID.IsZero() {
			j.ID = app.NewUlid()
		}
		j.User = user

		db := s3.NewWithContext(ctx)

		var b []byte
		if err = db.Put(j.Path(), app.MustMarshal(j)); err == nil {
			w := model.NewWork(j)
			j.WorkID = w.ID
			j.Err = db.Put(w.Path(), app.MustMarshal(w))
			b = app.MustMarshal(j)
			err = db.Put(j.Path(), b)
		}

		return app.Response(b, err)
	}

	return app.NotImplemented(req)
}
