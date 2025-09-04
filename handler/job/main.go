package main

import (
	"bytelyon-functions/internal/model"
	"bytelyon-functions/internal/util"
	"bytelyon-functions/pkg/api"
	"bytelyon-functions/pkg/service/s3"
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/oklog/ulid/v2"
)

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	api.LogURLRequest(req)

	if api.IsOptions(req) {
		return api.OK()
	}

	user, err := model.ValidateJWT(ctx, req.Headers["authorization"])
	if err != nil {
		return api.Unauthorized()
	}

	switch {
	case api.IsPatch(req):
		return handlePatch(ctx, user, req.QueryStringParameters["id"])
	case api.IsPost(req):
		return handleSave(ctx, user, req.Body, "")
	case api.IsGet(req):
		return handleGet(ctx, user, req.QueryStringParameters["size"], req.QueryStringParameters["after"])
	case api.IsPut(req):
		return handleSave(ctx, user, req.Body, req.QueryStringParameters["id"])
	case api.IsDelete(req):
		return handleDelete(ctx, user, req.QueryStringParameters["ids"])
	}

	return api.NotImplemented(req)
}

func handlePatch(ctx context.Context, u model.User, id string) (events.LambdaFunctionURLResponse, error) {

	ID, err := ulid.Parse(id)
	if err != nil {
		return api.BadRequest(err)
	}

	db := s3.NewWithContext(ctx)

	var out []byte
	if out, err = db.Get(ctx, "bytelyon", model.Job{ID: ID, User: u}.Key()); err == nil {
		var j model.Job
		util.MustUnmarshal(out, &j)
		w := model.NewWork(j)
		j.WorkID = w.ID
		j.Err = db.Put(ctx, "bytelyon", w.Key(), util.MustMarshal(w))
		err = db.Put(ctx, "bytelyon", j.Key(), util.MustMarshal(j))
	}

	return api.Err(err)
}

func handleSave(ctx context.Context, u model.User, body, id string) (events.LambdaFunctionURLResponse, error) {

	var j model.Job
	if err := json.Unmarshal([]byte(body), &j); err != nil {
		return api.BadRequest(err)
	} else if err = j.Validate(); err != nil {
		return api.BadRequest(err)
	}

	if id == "" {
		j.ID = model.NewUlid()
	} else {
		var err error
		if j.ID, err = ulid.Parse(id); err != nil {
			return api.BadRequest(err)
		}
	}
	j.User = u
	if err := s3.NewWithContext(ctx).Put(ctx, "bytelyon", j.Key(), util.MustMarshal(j)); err != nil {
		return api.ServerError(err)
	}

	return api.Marshall(j)
}

func handleGet(ctx context.Context, u model.User, size, after string) (events.LambdaFunctionURLResponse, error) {

	var n int32
	if i, err := strconv.Atoi(size); err != nil {
		n = 10
	} else {
		n = int32(i)
	}

	db := s3.NewWithContext(ctx)

	keys, err := db.KeysAfter(ctx, n, "bytelyon", model.Job{User: u}.Path(), after)
	if err != nil {
		return api.ServerError(err)
	}

	var vv []model.Job
	for _, key := range keys {

		var b []byte
		if b, err = db.Get(ctx, "bytelyon", key); err != nil {
			return api.ServerError(err)
		}

		var v model.Job
		util.MustUnmarshal(b, &v)
		vv = append(vv, v)
	}

	return api.Marshall(map[string]interface{}{
		"items": vv,
		"size":  len(vv),
	})
}

func handleDelete(ctx context.Context, u model.User, ids string) (events.LambdaFunctionURLResponse, error) {
	db := s3.NewWithContext(ctx)
	var err error
	for _, id := range strings.Split(ids, ",") {
		if ID, e := ulid.Parse(id); err != nil {
			err = errors.Join(err, e)
		} else if e = db.Delete(ctx, "bytelyon", model.Job{ID: ID, User: u}.Key()); e != nil {
			err = errors.Join(err, e)
		}
	}
	return api.Err(err)
}

func main() {
	lambda.Start(handler)
}
