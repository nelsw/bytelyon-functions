package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
		FormatLevel: func(a any) string {
			var color, level string
			switch level = strings.ToUpper(a.(string)[:3]); level {
			case "TRA":
				color = "\033[36m" // Cyan
			case "DEB":
				color = "\033[35m" // Magenta
			case "INF":
				color = "\033[32m" // Green
			case "WAR":
				color = "\033[33m" // Yellow
			case "ERR":
				color = "\033[31m" // Red
			case "FAT", "PAN":
				color = "\033[41m\033[37m" // Red background, white text
			default:
				color = "\033[0m" // Reset
			}
			return color + level + "\033[0m" // Reset color after level
		},
	})
}

type Arrangable interface {
	methodable
	Context(context.Context) Arrangable
	Header(string, string) Arrangable
	Param(string, string) Arrangable
	Body(any) Arrangable
}

func (h *helper) Context(c context.Context) Arrangable {
	h.ctx = c
	return h
}
func (h *helper) Header(k, v string) Arrangable {
	h.headers[k] = v
	return h
}
func (h *helper) Param(k, v string) Arrangable {
	h.params[k] = v
	return h
}
func (h *helper) Body(a any) Arrangable {
	_ = gofakeit.Struct(&a)
	b, _ := json.Marshal(&a)
	h.body = string(b)
	return h
}

type methodable interface {
	Method(string) Actable
	Delete(string) Actable
	Get() Actable
	Options() Actable
	Patch() Actable
	Post(any) Actable
	Put(any) Actable
}

func (h *helper) Method(s string) Actable {
	h.method = s
	return h
}
func (h *helper) Delete(s string) Actable {
	return h.Param("id", s).Method(http.MethodDelete)
}
func (h *helper) Get() Actable {
	return h.Method(http.MethodGet)
}
func (h *helper) Options() Actable {
	return h.Method(http.MethodOptions)
}
func (h *helper) Patch() Actable {
	return h.Method(http.MethodPatch)
}
func (h *helper) Post(a any) Actable {
	return h.Body(a).Method(http.MethodPost)
}
func (h *helper) Put(a any) Actable {
	return h.Body(a).Method(http.MethodPut)
}

type Handler = func(context.Context, events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error)
type Actable interface {
	Handle(func(context.Context, events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error)) Assertable
}

func (h *helper) Handle(f Handler) Assertable {
	h.res, _ = f(h.ctx, events.LambdaFunctionURLRequest{
		Headers: h.headers,
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Method: h.method,
			},
		},
		QueryStringParameters: h.params,
		Body:                  h.body,
	})
	return h
}

type Assertable interface {
	OK() Assertable
	Assert(string, any, any) Assertable
	JSON(map[string]any)
}

func (h *helper) OK() Assertable {
	return h.Assert("status", http.StatusOK, h.res.StatusCode)
}
func (h *helper) JSON(exp map[string]any) {

	var act map[string]any
	_ = json.Unmarshal([]byte(h.res.Body), &act)

	for k, v := range exp {
		val, ok := act[k]
		if !ok {
			val = nil
		}
		e := fmt.Sprintf("%v", v)
		a := fmt.Sprintf("%v", val)
		h.Assert(k, e, a)
	}

	b, _ := json.MarshalIndent(act, "", "\t")
	log.Debug().Msg(string(b))

	if h.failed {
		fmt.Printf("%+v\n", assert.CallerInfo()[1])
	}
}
func (h *helper) Assert(msg string, exp, act any) Assertable {

	level := zerolog.InfoLevel
	if ok := exp == act; !ok {
		level = zerolog.ErrorLevel
		h.failed = true
	}

	log.WithLevel(level).
		Any("got", act).
		Any("want", exp).
		Msgf(" %10s", msg)

	return h
}

type Helper interface {
	Arrangable
	Actable
	Assertable
}
type helper struct {
	tester
	requester
	responser
}
type tester struct {
	t      *testing.T
	failed bool
}
type requester struct {
	ctx     context.Context
	headers map[string]string
	params  map[string]string
	method  string
	body    string
}
type responser struct {
	res events.LambdaFunctionURLResponse
}

func New(t *testing.T) Helper {
	t.Setenv("APP_MODE", "local")
	return &helper{
		tester{t: t},
		requester{
			ctx:     context.Background(),
			headers: map[string]string{},
			params:  map[string]string{},
		},
		responser{},
	}
}
