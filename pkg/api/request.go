package api

import (
	"bytelyon-functions/pkg/model"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
)

type Request events.APIGatewayV2HTTPRequest

func NewRequest() *Request {
	return &Request{
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			Authorizer: &events.APIGatewayV2HTTPRequestContextAuthorizerDescription{
				Lambda: map[string]any{},
			},
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{},
		},
		QueryStringParameters: map[string]string{},
	}
}

func (r *Request) Log(s ...string) *Request {

	msgPrefix := "API"
	if len(s) > 0 {
		msgPrefix = s[0]
	}

	log.Info().
		Any("headers", r.Headers).
		Any("params", r.QueryStringParameters).
		Str("body", r.Body).
		Str("method", r.Method()).
		Str("user", r.User().ID.String()).
		Any(msgPrefix+" Request", r).
		Send()

	return r
}

func (r *Request) Method() string {
	return r.RequestContext.HTTP.Method
}

func (r *Request) User() *model.User {
	b, err := json.Marshal(r.RequestContext.Authorizer.Lambda["user"])
	if err != nil {
		log.Panic().Err(err).Send()
	}

	var u model.User
	if err = json.Unmarshal(b, &u); err != nil {
		log.Panic().Err(err).Send()
	}

	return &u
}

func (r *Request) Param(s string) string {
	return r.QueryStringParameters[s]
}

func (r *Request) Data() []byte {
	return []byte(r.Body)
}

func (r *Request) WithUser(v model.User) *Request {
	r.RequestContext.Authorizer.Lambda["user"] = v
	return r
}

func (r *Request) WithParam(k, v string) *Request {
	r.QueryStringParameters[k] = v
	return r
}

func (r *Request) WithData(a any) *Request {
	b, _ := json.Marshal(a)
	r.Body = string(b)
	return r
}

func (r *Request) Delete() Request {
	return *r.method(http.MethodDelete)
}

func (r *Request) Get() Request {
	return *r.method(http.MethodGet)
}

func (r *Request) Post() Request {
	return *r.method(http.MethodPost)
}

func (r *Request) Put() Request {
	return *r.method(http.MethodPut)
}

func (r *Request) method(s string) *Request {
	r.RequestContext.HTTP.Method = s
	return r
}
