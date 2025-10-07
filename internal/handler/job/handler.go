package job

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/s3"
	"bytelyon-functions/internal/handler/jwt"
	"bytelyon-functions/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/oklog/ulid/v2"
)

func Handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	app.LogURLRequest(req)

	if app.IsOptions(req) {
		return app.OK()
	}

	user, err := jwt.Validate(ctx, req.Headers["authorization"])
	if err != nil {
		return app.Unauthorized(err)
	}

	db := s3.New(ctx)

	if app.IsDelete(req) {
		return app.Err(Delete(db, user.ID, req.QueryStringParameters["id"]))
	}

	if app.IsPut(req) || app.IsPost(req) {
		var job model.Job
		if job, err = Save(db, user.ID, []byte(req.Body)); err != nil {
			return app.BadRequest(err)
		}
		return app.Marshall(job)
	}

	if app.IsGet(req) {
		i, _ := strconv.Atoi(req.QueryStringParameters["size"])
		if i < 1 || i > 10 {
			i = 10
		}
		var page model.Page
		if page, err = FindAll(db, user, i, req.QueryStringParameters["after"]); err != nil {
			return app.BadRequest(err)
		}
		return app.Marshall(page)
	}

	return app.NotImplemented(req)
}

func Save(db s3.Client, userID ulid.ULID, in []byte) (job model.Job, err error) {

	err = json.Unmarshal(in, &job)
	if err != nil {
		return
	}
	if err = job.Validate(); err != nil {
		return
	}

	var run bool
	if job.ID.IsZero() {
		job.ID = app.NewUlid()
		run = true
	}

	if job.Type == model.NewsJobType {
		var keywordQuery string
		for i, keyword := range job.Keywords {
			if i > 0 {
				keywordQuery += "+"
			}
			keywordQuery += url.QueryEscape(keyword)
		}
		job.URLs = []string{
			fmt.Sprintf("https://news.google.com/rss/search?q=%s&hl=en-US&gl=US&ceid=US:en", keywordQuery),
			fmt.Sprintf("https://www.bing.com/news/search?format=rss&q=%s", keywordQuery),
			fmt.Sprintf("https://www.bing.com/search?format=rss&q=%s", keywordQuery),
		}
	}

	if err = db.Put(model.JobKey(userID, job.ID), app.MustMarshal(job)); err != nil {
		return
	}

	if run {
		job.DoWork(db, userID)
	}

	return
}

func FindAll(db s3.Client, user model.User, size int, after string) (page model.Page, err error) {
	return model.FindJobsFast(db, user.ID)
}

func Delete(db s3.Client, userID ulid.ULID, id string) error {
	return db.Delete(model.JobKey(userID, id))
}
