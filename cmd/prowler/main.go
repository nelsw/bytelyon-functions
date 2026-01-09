package main

import (
	"bytelyon-functions/pkg/api"
	"bytelyon-functions/pkg/db"
	"bytelyon-functions/pkg/model"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/oklog/ulid/v2"
)

func Handler(r api.Request) (events.APIGatewayV2HTTPResponse, error) {
	r.Log()
	switch r.Method() {
	case http.MethodPut:
		return handlePut(r)
	case http.MethodGet:
		return handleGet(r)
	case http.MethodDelete:
		return handleDelete(r)
	}
	return api.NotImplemented()
}

func handlePut(r api.Request) (events.APIGatewayV2HTTPResponse, error) {

	var p = new(model.Prowler)
	if err := json.Unmarshal([]byte(r.Body), p); err != nil {
		return api.BadRequest(err)
	} else if err = p.Type.Validate(); err != nil {
		return api.BadRequest(err)
	}

	if p.ID == "" {
		return api.BadRequest(errors.New("id required"))
	}

	if strings.HasSuffix(p.ID, "/") {
		p.ID = strings.TrimSuffix(p.ID, "/")
	}

	p.UserID = r.User().ID

	if err := db.Save(p); err != nil {
		return api.BadRequest(err)
	}

	return api.Response(p, db.Find(p))
}

// todo - delete ids
func handleDelete(r api.Request) (events.APIGatewayV2HTTPResponse, error) {
	id, err := ulid.Parse(r.Param("id"))
	if err != nil {
		return api.BadRequest(err)
	} else if err = db.MagicDelete(r.User().ID, id); err != nil {
		return api.BadRequest(err)
	}
	return api.OK()
}

func handleGet(r api.Request) (events.APIGatewayV2HTTPResponse, error) {

	t, err := model.NewProwlerType(r.Param("type"))
	if err != nil {
		return api.BadRequest(err)
	}

	p := &model.Prowler{
		UserID: r.User().ID,
		Type:   t,
		ID:     r.Param("id"),
	}

	return api.Response(p.FindAll())
}

func main() {
	lambda.Start(Handler)
}
