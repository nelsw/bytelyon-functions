package main

import (
	"bytelyon-functions/internal/entity"
	"bytelyon-functions/internal/model"
	"bytelyon-functions/pkg/api"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	api.LogURLRequest(req)
	switch req.RequestContext.HTTP.Method {
	case http.MethodOptions:
		return api.OK()
	case http.MethodPatch:
		return handlePatch(ctx, req.QueryStringParameters["id"])
	case http.MethodPost:
		return handleSave(ctx, req.Body, model.Job{ID: model.NewUlid()})
	case http.MethodGet:
		return handleGet(ctx, req.QueryStringParameters["size"])
	case http.MethodPut:
		return handleSave(ctx, req.Body, model.Job{})
	case http.MethodDelete:
		return handleDelete(ctx, req.QueryStringParameters["ids"])
	default:
		return api.NotImplemented(req.RequestContext.HTTP.Method)
	}
}

func handlePatch(ctx context.Context, id string) (events.LambdaFunctionURLResponse, error) {
	var v model.Job
	if err := entity.New(ctx).Value(&v).ID(id).Find(); err != nil {
		return api.BadRequest(err)
	} else if err = v.CreateWork(); err != nil {
		return api.ServerError(err)
	}
	return api.OK()
}

func handleSave(ctx context.Context, body string, j model.Job) (events.LambdaFunctionURLResponse, error) {
	if err := json.Unmarshal([]byte(body), &j); err != nil {
		return api.BadRequest(err)
	} else if err = j.Validate(); err != nil {
		return api.BadRequest(err)
	} else if err = entity.New(ctx).Value(&j).Save(); err != nil {
		return api.ServerError(err)
	}
	return api.OK(&j)
}

func handleGet(ctx context.Context, size string) (events.LambdaFunctionURLResponse, error) {

	n, err := strconv.Atoi(size)
	if err != nil {
		n = 10
	}

	var vv []model.Job
	if err = entity.New(ctx).Value(model.Job{}).Type(&vv).Page(int32(n)); err != nil {
		return api.ServerError(err)
	}

	return api.OK(map[string]interface{}{
		"items": vv,
		"size":  len(vv),
	})
}

func handleDelete(ctx context.Context, ids string) (events.LambdaFunctionURLResponse, error) {
	var err error
	for _, v := range strings.Split(ids, ",") {
		if e := entity.New(ctx).Value(&model.Job{}).ID(v).Delete(); e != nil {
			err = errors.Join(err, e)
		}
	}
	if err != nil {
		return api.BadRequest(err)
	}
	return api.OK()
}

func main() {
	lambda.Start(handler)
}
