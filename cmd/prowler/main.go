package main

import (
	"bytelyon-functions/pkg/api"
	"bytelyon-functions/pkg/db"
	"bytelyon-functions/pkg/model"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/oklog/ulid/v2"
)

func Handler(r api.Request) (events.APIGatewayV2HTTPResponse, error) {
	r.Log()
	switch r.Method() {
	case http.MethodPatch:
		return handlePatch(r)
	case http.MethodPost:
		return handlePost(r)
	case http.MethodGet:
		return handleGet(r)
	}
	return api.NotImplemented()
}

func handlePatch(r api.Request) (events.APIGatewayV2HTTPResponse, error) {

	var err error
	p := &model.Prowler{
		UserID: r.User().ID,
	}

	if p.Type, err = model.NewProwlerType(r.Param("type")); err != nil {
		return api.BadRequest(err)
	} else if p.ID, err = ulid.Parse(r.Param("id")); err != nil {
		return api.BadRequest(err)
	}

	var data struct {
		model.Targets `json:"targets"`
		time.Duration `json:"frequency"`
	}

	if err = json.Unmarshal([]byte(r.Body), &data); err != nil {
		return api.BadRequest(err)
	}

	p.Targets = data.Targets
	p.Frequency = data.Duration

	return api.Response(p, db.Save(p))
}

func handlePost(r api.Request) (events.APIGatewayV2HTTPResponse, error) {

	var p = new(model.Prowler)
	if err := json.Unmarshal([]byte(r.Body), p); err != nil {
		return api.BadRequest(err)
	} else if err = p.Type.Validate(); err != nil {
		return api.BadRequest(err)
	}

	if p.Frequency < time.Duration(10)*time.Minute {
		return api.BadRequest(errors.New("sitemap prowl frequency must be at least 10 minutes"))
	} else if p.Type == model.SitemapProwlerType && !strings.HasPrefix(p.URL, "https://") {
		return api.BadRequest(errors.New("sitemap prowl url must be set"))
	} else if p.Query == "" {
		return api.BadRequest(errors.New("query must be set"))
	}
	// todo - improve type specific validation

	if strings.HasSuffix(p.URL, "/") {
		p.URL = strings.TrimSuffix(p.URL, "/")
	}

	p.UserID = r.User().ID
	p.ID = model.NewUlid()

	if err := db.Save(p); err != nil {
		return api.BadRequest(err)
	}

	p.Prowl()

	return api.Response(p, db.Find(p))
}

func handleGet(r api.Request) (events.APIGatewayV2HTTPResponse, error) {

	t, err := model.NewProwlerType(r.Param("type"))
	if err != nil {
		return api.BadRequest(err)
	}

	p := &model.Prowler{
		UserID: r.User().ID,
		Type:   t,
	}

	// find all
	if r.Param("id") == "" {
		return api.Response(db.List(p))
	}

	// find one
	if p.ID, err = ulid.Parse(r.Param("id")); err != nil {
		return api.BadRequest(errors.Join(err, errors.New("invalid id")))
	} else if err = db.Find(p); err != nil {
		return api.BadRequest(errors.Join(err, errors.New("not found")))
	}
	return api.OK(p)
}

func main() {
	lambda.Start(Handler)
}
