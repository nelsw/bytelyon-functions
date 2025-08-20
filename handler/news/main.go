package main

import (
	"bytelyon-functions/internal/entity"
	"bytelyon-functions/internal/model/news"
	"bytelyon-functions/pkg/api"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	api.LogURLRequest(req)

	switch req.RequestContext.HTTP.Method {
	case http.MethodOptions:
		return api.Response(http.StatusOK, "")
	case http.MethodGet:
		return handleGet(ctx,
			req.QueryStringParameters["type"],
			req.QueryStringParameters["id"],
			req.QueryStringParameters["size"])
	case http.MethodPatch:
		return handlePatch(ctx, req.QueryStringParameters["id"])
	case http.MethodPut:
		return handlePut(ctx, req.Body)
	case http.MethodPost:
		return handlePost(ctx, req.Body)
	case http.MethodDelete:
		return handleDelete(ctx, req.QueryStringParameters["type"], req.QueryStringParameters["id"])
	default:
		return api.Response(http.StatusNotImplemented, "")
	}
}

func handleGet(ctx context.Context, t, id, size string) (events.LambdaFunctionURLResponse, error) {

	if id != "" {
		switch t {
		case "job":
			var v news.Job
			if err := entity.New(ctx).Value(&v).Find(); err != nil {
				return api.Response(http.StatusBadRequest, err.Error())
			}
			b, _ := json.Marshal(&v)
			return api.Response(http.StatusOK, string(b))
		case "work":
			var v news.Work
			if err := entity.New(ctx).Value(&v).Find(); err != nil {
				return api.Response(http.StatusBadRequest, err.Error())
			}
			b, _ := json.Marshal(&v)
			return api.Response(http.StatusOK, string(b))
		case "article":
			fallthrough
		default:
			return api.Response(http.StatusNotImplemented, "")
		}
	}

	n, err := strconv.Atoi(size)
	if err != nil {
		return api.Response(http.StatusBadRequest, err.Error())
	}

	switch t {
	case "job":
		var vv []news.Job
		if err = entity.New(ctx).Value(news.Job{}).Type(&vv).Page(int32(n)); err != nil {
			return api.Response(http.StatusInternalServerError, err.Error())
		} else {
			b, _ := json.Marshal(&vv)
			return api.Response(http.StatusOK, string(b))
		}
	case "work":
		var vv []news.Work
		if err = entity.New(ctx).Value(news.Work{}).Type(&vv).Page(int32(n)); err != nil {
			return api.Response(http.StatusInternalServerError, err.Error())
		} else {
			b, _ := json.Marshal(&vv)
			return api.Response(http.StatusOK, string(b))
		}
	default:
		return api.Response(http.StatusNotImplemented, t)
	}
}

func handlePut(ctx context.Context, body string) (events.LambdaFunctionURLResponse, error) {
	var v news.Job
	if err := json.Unmarshal([]byte(body), &v); err != nil {
		return api.Response(http.StatusBadRequest, err.Error())
	}
	v.UpdatedAt = time.Now()
	if err := entity.New(ctx).Value(&v).Save(); err != nil {
		return api.Response(http.StatusInternalServerError, err.Error())
	}
	b, _ := json.Marshal(&v)
	return api.Response(http.StatusOK, string(b))
}

func handlePatch(ctx context.Context, id string) (events.LambdaFunctionURLResponse, error) {
	var v news.Job
	if err := entity.New(ctx).Value(&v).Find(); err != nil {
		return api.Response(http.StatusBadRequest, err.Error())
	}

	if err := v.CreateWork(); err != nil {
		return api.Response(http.StatusInternalServerError, err.Error())
	}

	return api.Response(http.StatusOK, "")
}

func handlePost(ctx context.Context, body string) (events.LambdaFunctionURLResponse, error) {

	if body == "" {
		return api.Response(http.StatusBadRequest, "must provide id or body")
	}

	var v news.Job
	if err := json.Unmarshal([]byte(body), &v); err != nil {
		return api.Response(http.StatusBadRequest, err.Error())
	}

	v.ID = uuid.New()
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()

	if err := entity.New(ctx).Value(&v).Save(); err != nil {
		return api.Response(http.StatusInternalServerError, err.Error())
	}

	b, _ := json.Marshal(&v)
	return api.Response(http.StatusOK, string(b))
}

func handleDelete(ctx context.Context, t string, id string) (events.LambdaFunctionURLResponse, error) {
	ids := strings.Split(id, ",")
	switch t {
	case "job":
		var v news.Job
		for _, s := range ids {
			if err := entity.New(ctx).Value(&v).ID(s).Delete(); err != nil {
				return events.LambdaFunctionURLResponse{}, err
			}
		}
		return api.Response(http.StatusOK, "")
	case "work":
		var v news.Work
		for _, s := range ids {
			if err := entity.New(ctx).Value(&v).ID(s).Delete(); err != nil {
				return events.LambdaFunctionURLResponse{}, err
			}
		}
		return api.Response(http.StatusOK, "")
	case "article":
		fallthrough
	default:
		return api.Response(http.StatusNotImplemented, "")
	}
}

func main() {
	lambda.Start(handler)
}
