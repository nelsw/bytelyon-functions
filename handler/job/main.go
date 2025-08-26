package main

import (
	"bytelyon-functions/internal/entity"
	"bytelyon-functions/internal/model/bot"
	"bytelyon-functions/internal/model/id"
	"bytelyon-functions/pkg/api"
	"context"
	"encoding/json"
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
		return handlePost(ctx, req.Body)
	case http.MethodGet:
		return handleGet(ctx, req.QueryStringParameters["size"])
	case http.MethodPut:
		return handlePut(ctx, req.Body)
	case http.MethodDelete:
		return handleDelete(ctx, req.QueryStringParameters["ids"])
	default:
		return api.NotImplemented(req.RequestContext.HTTP.Method)
	}
}

func handlePatch(ctx context.Context, id string) (events.LambdaFunctionURLResponse, error) {
	var v bot.Job
	if err := entity.New(ctx).Value(&v).ID(id).Find(); err != nil {
		return api.BadRequest(err)
	} else if err = v.CreateWork(); err != nil {
		return api.ServerError(err)
	}
	return api.OK()
}

func handlePost(ctx context.Context, body string) (events.LambdaFunctionURLResponse, error) {

	var v bot.Job
	if err := json.Unmarshal([]byte(body), &v); err != nil {
		return api.BadRequest(err)
	} else if err = v.Validate(); err != nil {
		return api.BadRequest(err)
	}

	if v.ID = id.NewULID(); v.Name == "" {
		v.Name = v.ID.String()
	}
	v.CreatedAt = bot.NewTimeStamp()
	v.UpdatedAt = bot.NewTimeStamp()

	if err := entity.New(ctx).Value(&v).Save(); err != nil {
		return api.ServerError(err)
	}

	return api.OK(&v)
}

func handlePut(ctx context.Context, body string) (events.LambdaFunctionURLResponse, error) {
	var v bot.Job
	if err := json.Unmarshal([]byte(body), &v); err != nil {
		return api.BadRequest(err)
	} else if err = v.Validate(); err != nil {
		return api.BadRequest(err)
	}
	v.UpdatedAt = bot.NewTimeStamp()
	if err := entity.New(ctx).Value(&v).Save(); err != nil {
		return api.ServerError(err)
	}
	return api.OK(&v)
}

func handleGet(ctx context.Context, size string) (events.LambdaFunctionURLResponse, error) {

	n, err := strconv.Atoi(size)
	if err != nil {
		n = 10
	}

	var vv []bot.Job
	if err = entity.New(ctx).Value(bot.Job{}).Type(&vv).Page(int32(n)); err != nil {
		return api.ServerError(err)
	}

	return api.OK(map[string]interface{}{
		"items": vv,
		"size":  len(vv),
	})
}

func handleDelete(ctx context.Context, ids string) (events.LambdaFunctionURLResponse, error) {
	var j bot.Job
	m := map[string]string{}
	for _, v := range strings.Split(ids, ",") {
		if err := entity.New(ctx).Value(&j).ID(v).Delete(); err != nil {
			m[v] = err.Error()
		}
	}
	return api.OK(map[string]interface{}{"errors": m})
}

func main() {
	lambda.Start(handler)
}
