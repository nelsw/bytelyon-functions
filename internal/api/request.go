package api

import (
	"bytelyon-functions/internal/model"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
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

func (r *Request) User(v model.User) *Request {
	r.RequestContext.Authorizer.Lambda["user"] = v
	return r
}

func (r *Request) Method(s string) *Request {
	r.RequestContext.HTTP.Method = s
	return r
}

func (r *Request) Param(k, v string) *Request {
	r.QueryStringParameters[k] = v
	return r
}

func (r *Request) Data(a any) *Request {
	b, _ := json.Marshal(a)
	r.Body = string(b)
	return r
}

func (r *Request) Build() events.APIGatewayV2HTTPRequest {
	return events.APIGatewayV2HTTPRequest(*r)
}
