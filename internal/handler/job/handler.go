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
		if job, err = Save(db, user.ID, []byte(req.Body), app.IsPost(req)); err != nil {
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

func Save(db s3.Client, userID ulid.ULID, in []byte, run bool) (job model.Job, err error) {

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
			"https://news.google.com/rss/search?q={KEYWORD_QUERY}&hl=en-US&gl=US&ceid=US:en",
			"https://www.bing.com/news/search?format=rss&q={KEYWORD_QUERY}",
			"https://www.bing.com/search?format=rss&q={KEYWORD_QUERY}",
		}
	}

	if err = db.Put(model.JobKey(userID, job.ID), app.MustMarshal(job)); err != nil {
		return
	}

	if run {
		work.Now(db, job)
	}

	return
}

func FindAll(db s3.Client, user model.User, size int, after string) (page model.Page, err error) {

	var jobs model.Jobs
	if jobs, err = user.FindAllJobs(db); err != nil {
		return
	} else if page.Total = len(jobs); page.IsEmpty() {
		return
	}

	afterFound := after == ""
	for _, job := range jobs {

		if page.Size >= size {
			break
		}

		if !afterFound && job.ID.String() == after {
			afterFound = true
		}

		if !afterFound {
			continue
		}

		var w model.Work
		if e := db.Find(model.JobKey(user.ID, job.ID), &w.Job); e != nil {
			err = errors.Join(err, e)
		} else if w.Items, e = w.Job.Items(db); e != nil {
			err = errors.Join(err, e)
		} else {
			page = page.AddItem(w)
		}
	}

	return
}

func Delete(db s3.Client, userID ulid.ULID, id string) error {
	return db.Delete(model.JobKey(userID, id))
}
