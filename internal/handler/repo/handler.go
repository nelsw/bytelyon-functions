package repo

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/s3"
	"bytelyon-functions/internal/model"
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/oklog/ulid/v2"
)

func Handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	app.LogURLRequest(req)

	if app.IsOptions(req) {
		return app.OK()
	}

	user, err := model.ValidateJWT(ctx, req.Headers["authorization"])
	if err != nil {
		return app.Unauthorized(err)
	} else if !strings.Contains("user/"+user.ID.String(), req.RawPath) {
		return app.Forbidden()
	}

	path := req.RawPath
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	db := s3.NewWithContext(ctx)
	if app.IsDelete(req) {
		return app.Err(db.Delete(path))
	}

	if app.IsGet(req) {
		parts := strings.Split(path, "/")
		lastPart := parts[len(parts)-1]
		if _, err = ulid.Parse(lastPart); err != nil {
			return app.Response(db.Get(path))
		}

		n, _ := strconv.Atoi(req.QueryStringParameters["size"])
		var keys []string
		if keys, err = db.Keys(path, req.QueryStringParameters["after"], n); err != nil {
			return app.BadRequest(err)
		}

		var ss []string
		for _, key := range keys {
			if o, e := db.Get(key); e != nil {
				err = errors.Join(err, e)
			} else {
				ss = append(ss, string(o))
			}
		}

		return app.Marshall(map[string]interface{}{
			"items": ss,
			"size":  len(ss),
		})
	}

	return app.NotImplemented(req)
}
