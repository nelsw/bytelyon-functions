package login

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/s3"
	"bytelyon-functions/internal/model"
	"context"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func Handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	app.LogURLRequest(req)

	if app.IsOptions(req) {
		return app.OK()
	}

	if app.IsPost(req) && req.RawPath == "/login" {
		token := req.Headers["authorization"]
		if strings.HasPrefix(token, "Basic ") {
			token = strings.TrimPrefix(token, "Basic ")
		}

		b, _ := base64.StdEncoding.DecodeString(token)
		parts := strings.Split(string(b), ":")
		if len(parts) != 2 {
			return app.BadRequest(errors.New("invalid token"))
		}

		db := s3.NewWithContext(ctx)

		email := model.Email{ID: parts[0]}
		if out, err := db.Get(email.Path()); err != nil {
			return app.BadRequest(err)
		} else {
			app.MustUnmarshal(out, &email)
		}

		pass := model.Password{User: model.User{ID: email.UserID}, Text: parts[1]}
		if out, err := db.Get(pass.Path()); err != nil {
			return app.BadRequest(err)
		} else {
			app.MustUnmarshal(out, &pass)
			if pass.Compare() != nil {
				return app.Unauthorized(errors.New("invalid password"))
			}
		}

		return app.Response(model.CreateJWT(ctx, model.User{ID: email.UserID}))
	}

	return app.NotImplemented(req)
}
