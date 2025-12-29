package main

import (
	"bytelyon-functions/pkg/api"
	"bytelyon-functions/pkg/db"
	"bytelyon-functions/pkg/model"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(req api.Request) (events.APIGatewayV2HTTPResponse, error) {

	req.Log()

	p := &model.Profile{UserID: req.User().ID}

	switch req.Method() {
	case http.MethodGet:
		return api.Response(p, db.Find(p))
	case http.MethodPut:
		if err := json.Unmarshal([]byte(req.Body), p); err != nil {
			return api.BadRequest(err)
		}
		p.UserID = req.User().ID
		return api.Response(p, db.Save(p))
	}

	return api.NotImplemented()
}

func main() {
	lambda.Start(Handler)
}
