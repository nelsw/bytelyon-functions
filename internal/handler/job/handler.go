package job

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/s3"
	"bytelyon-functions/internal/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

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
	}

	db := s3.NewWithContext(ctx)

	if app.IsDelete(req) {
		return app.Err(Delete(db, user, req.QueryStringParameters["id"]))
	}

	if app.IsPut(req) || app.IsPost(req) {
		return app.Response(Save(db, user, req.Body, app.IsPost(req)))
	}

	if app.IsGet(req) {
		i, _ := strconv.Atoi(req.QueryStringParameters["size"])
		return app.Response(FindAll(db, user, i, req.QueryStringParameters["after"]))
	}

	return app.NotImplemented(req)
}

func Save(db s3.Client, u model.User, s string, run bool) (b []byte, err error) {
	var j model.Job

	if j, err = model.MakeJob(u, []byte(s)); err == nil {
		if run {
			w := model.NewWork(j)
			j.WorkID = w.ID
			if err = db.Put(w.Key(), app.MustMarshal(w)); err != nil {
				err = errors.Join(err, errors.New("failed to put work"))
			}
		}
		b = app.MustMarshal(j)
		err = db.Put(j.Key(), b)
	}

	return
}

func FindAll(db s3.Client, u model.User, size int, after string) ([]byte, error) {
	if size == 0 {
		size = 10
	}
	keys, err := db.Keys(model.Job{User: u}.Path(), after, 1000)
	var jj []model.Job
	if err == nil {
		for i, key := range keys {
			if o, e := db.Get(key); e != nil {
				err = errors.Join(err, e)
			} else {
				var j model.Job
				j.User = u
				app.MustUnmarshal(o, &j)

				kk, ee := db.Keys(model.MakeWork(j).Path(), "", 1000)
				fmt.Println(kk, ee)
				if ee != nil {
					err = errors.Join(err, ee)
				} else {
					for _, k := range kk {
						if o, e = db.Get(k); e != nil {
							err = errors.Join(err, e)
						} else {
							var w model.Work
							if e = json.Unmarshal(o, &w); e != nil {
								err = errors.Join(err, e)
							} else {
								j.Work = append(j.Work, w)
							}
						}
					}
				}

				jj = append(jj, j)
			}
			if i >= size {
				break
			}
		}
	}

	return app.MustMarshal(map[string]any{"items": jj, "size": len(jj)}), err
}

func Delete(db s3.Client, u model.User, id string) error {
	ID, err := ulid.Parse(id)
	if err == nil {
		err = db.Delete(model.Job{User: u, ID: ID}.Key())
	}
	return err
}
