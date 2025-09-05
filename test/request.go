package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

type Methodable interface {
	Method(string) Builder
	Options() events.LambdaFunctionURLRequest
	Post(...any) events.LambdaFunctionURLRequest
	Patch() events.LambdaFunctionURLRequest
	Get() events.LambdaFunctionURLRequest
	Delete() events.LambdaFunctionURLRequest
	Put() events.LambdaFunctionURLRequest
}

type Builder interface {
	Methodable
	Headers(map[string]string) Builder
	Header(string, string) Builder
	Bearer(string) Builder
	Path(string) Builder
	Query(string, any) Builder
	Body(any) Builder
	Build() events.LambdaFunctionURLRequest
}

type builder struct {
	method  string
	headers map[string]string
	path    []string
	query   map[string]string
	body    string
}

func (b builder) Body(a any) Builder {
	if a == nil {
		return b
	}
	if s, ok := a.(string); ok {
		b.body = s
	} else {
		out, _ := json.Marshal(&a)
		b.body = string(out)
	}
	return b
}

func (b builder) Path(s string) Builder {
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}
	b.path = append(b.path, s)
	return b
}

func (b builder) Query(k string, v any) Builder {
	b.query[k] = fmt.Sprintf("%v", v)
	return b
}

func (b builder) Headers(m map[string]string) Builder {
	b.headers = m
	return b
}

func (b builder) Header(k, v string) Builder {
	b.headers[k] = v
	return b
}

func (b builder) Bearer(v string) Builder {
	b.headers["authorization"] = "Bearer " + v
	return b
}

func (b builder) Method(method string) Builder {
	b.method = method
	return b
}

func (b builder) Options() events.LambdaFunctionURLRequest {
	return b.Method(http.MethodOptions).Build()
}

func (b builder) Post(aa ...any) events.LambdaFunctionURLRequest {
	var a any
	if aa != nil && len(aa) == 1 && aa[0] != nil {
		a = aa[0]
	}
	return b.Body(a).Method(http.MethodPost).Build()
}

func (b builder) Patch() events.LambdaFunctionURLRequest {
	return b.Method(http.MethodPatch).Build()
}

func (b builder) Get() events.LambdaFunctionURLRequest {
	return b.Method(http.MethodGet).Build()
}

func (b builder) Delete() events.LambdaFunctionURLRequest {
	return b.Method(http.MethodDelete).Build()
}

func (b builder) Put() events.LambdaFunctionURLRequest {
	return b.Method(http.MethodPut).Build()
}

func (b builder) Build() events.LambdaFunctionURLRequest {
	req := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: b.method,
			},
		},
	}
	if len(b.headers) > 0 {
		req.Headers = b.headers
	}
	if len(b.path) > 0 {
		req.RawPath = strings.Join(b.path, "")
	}
	if len(b.query) > 0 {
		req.QueryStringParameters = b.query
		var params []string
		for k, v := range b.query {
			params = append(params, fmt.Sprintf("%s=%s", k, v))
		}
		req.RawQueryString = "?" + strings.Join(params, "&")
	}
	if len(b.body) > 0 {
		req.Body = b.body
	}
	return req
}

func NewRequest(t *testing.T) Builder {
	InitLogger()
	t.Setenv("APP_MODE", "test")
	t.Setenv("JWT_SECRET", "a-string-secret-at-least-256-bits-long")
	return builder{
		headers: map[string]string{},
		path:    []string{},
		query:   map[string]string{},
	}
}
