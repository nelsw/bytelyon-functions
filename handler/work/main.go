package main

import (
	"bytelyon-functions/internal/entity"
	"bytelyon-functions/internal/model"
	"bytelyon-functions/pkg/api"
	"context"
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
	case http.MethodGet:
		return handleGet(ctx, req.QueryStringParameters["size"])
	case http.MethodDelete:
		return handleDelete(ctx, req.QueryStringParameters["ids"])
	default:
		return api.NotImplemented(req.RequestContext.HTTP.Method)
	}
}

func handleGet(ctx context.Context, size string) (events.LambdaFunctionURLResponse, error) {

	n, err := strconv.Atoi(size)
	if err != nil {
		n = 10
	}

	var vv []model.Work
	if err = entity.New(ctx).Value(model.Work{}).Type(&vv).Page(int32(n)); err != nil {
		return api.ServerError(err)
	}
	return api.OK(&vv)
}

func handleDelete(ctx context.Context, ids string) (events.LambdaFunctionURLResponse, error) {
	var v model.Work
	m := map[string]string{}
	for _, id := range strings.Split(ids, ",") {
		if err := entity.New(ctx).Value(&v).ID(id).Delete(); err != nil {
			m[id] = err.Error()
		}
	}
	return api.OK(map[string]interface{}{"errors": m})
}

func main() {
	lambda.Start(handler)
}
