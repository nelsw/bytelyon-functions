package job

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/s3"
	"bytelyon-functions/internal/handler/jwt"
	"bytelyon-functions/internal/handler/work"
	"bytelyon-functions/internal/model"
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/oklog/ulid/v2"
)

func Handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	app.LogURLRequest(req)

	if app.IsOptions(req) {
		return app.OK()
	}

	var userID ulid.ULID
	if user, err := jwt.Validate(ctx, req.Headers["authorization"]); err != nil {
		return app.Unauthorized(err)
	} else {
		userID = user.ID
	}

	if app.IsDelete(req) {
		return app.Err(Delete(ctx, userID, req.QueryStringParameters["id"]))
	}

	if app.IsPut(req) || app.IsPost(req) {
		job, err := Save(ctx, userID, []byte(req.Body), app.IsPost(req))
		if err != nil {
			return app.BadRequest(err)
		}
		return app.Marshall(job)
	}

	if app.IsGet(req) {
		i, _ := strconv.Atoi(req.QueryStringParameters["size"])
		page, err := FindAll(ctx, userID, i, req.QueryStringParameters["after"])
		if err != nil {
			return app.BadRequest(err)
		}
		return app.Marshall(page)
	}

	return app.NotImplemented(req)
}

func Save(ctx context.Context, userID ulid.ULID, in []byte, run bool) (job model.Job, err error) {

	_ = json.Unmarshal(in, &job)
	if err = job.Validate(); err != nil {
		return
	}

	job.UserID = userID

	if job.ID.IsZero() {
		job.ID = app.NewUlid()
	}

	if job.Type == model.NewsJobType {
		job.URLs = []string{
			"https://news.google.com/rss/search?q=%s&hl=en-US&gl=US&ceid=US:en",
			"https://www.bing.com/news/search?format=rss&q=%s",
			"https://www.bing.com/search?format=rss&q=%s",
		}
	}

	if err = s3.New(ctx).Put(job.Key(), app.MustMarshal(job)); err != nil {
		return
	}

	if run {
		err = work.Job(job)
	}

	return
}

func FindAll(ctx context.Context, userID ulid.ULID, size int, after string) (page model.Page, err error) {

	if size < 1 || size > 10 {
		size = 10
	}

	db := s3.New(ctx)
	prefix := model.Job{UserID: userID}.Path()

	for {
		var keys []string
		if keys, err = db.Keys(prefix, after, 1000); err != nil {
			return
		}

		if len(keys) == 0 {
			return
		}

		page.Total += len(keys)

		if page.Items == nil {
			for _, key := range keys {

				o, e := db.Get(key)
				if e != nil {
					err = errors.Join(err, e)
					continue
				}

				var job model.Job
				app.MustUnmarshal(o, &job)
				// todo - items
				page.Items = append(page.Items, job)
				page.Size += 1
			}
		}

		if len(keys) == 1000 {
			after = keys[999]
			continue
		}

		return
	}
}

func Delete(ctx context.Context, userID ulid.ULID, id string) error {
	ID, err := ulid.Parse(id)
	if err != nil {
		return err
	}
	return s3.New(ctx).Delete(model.Job{ID: ID, UserID: userID}.Key())
}
