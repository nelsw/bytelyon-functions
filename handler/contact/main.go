package main

import (
	"bytelyon-functions/pkg/api"
	"bytelyon-functions/pkg/service/s3"
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/oklog/ulid/v2"
)

type Contact struct {
	ID    ulid.ULID `json:"id" fake:"skip"`
	Name  string    `json:"name" fake:"{name}"`
	Email string    `json:"email" fake:"{email}"`
	Value string    `json:"message" fake:"{sentence}"`
}

const bucket = "bytelyon"

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	api.LogURLRequest(req)

	var b []byte
	var err error

	switch {
	case api.IsOptions(req):
		return api.OK()
	case api.IsPost(req):
		b, err = handlePost(ctx, req.Body)
	case api.IsGet(req):
		b, err = handleGet(ctx)
	case api.IsPatch(req):
		err = handlePatch(ctx, req.QueryStringParameters["id"], req.QueryStringParameters["read"])
	case api.IsDelete(req):
		b, err = handleDelete(ctx, req.QueryStringParameters["ids"])
	default:
		return api.NotImplemented(req)
	}

	return api.Response(b, err)
}

func handlePost(ctx context.Context, s string) ([]byte, error) {

	var c Contact

	// unmarshal & validate
	if err := json.Unmarshal([]byte(s), &c); err != nil {
		return nil, err
	} else if c.Name == "" {
		return nil, errors.New("name is required")
	} else if c.Email == "" {
		return nil, errors.New("email is required")
	} else if c.Value == "" {
		return nil, errors.New("message is required")
	}

	// make ID
	c.ID = ulid.Make()

	// define key
	key := "message/contact/unread/" + c.ID.String() + ".json"

	// marshall
	b, _ := json.Marshal(c)

	return b, s3.NewClient(ctx).Put(ctx, bucket, key, b)
}

func handleGet(ctx context.Context) ([]byte, error) {
	out, err := s3.NewClient(ctx).List(ctx, bucket, "message/contact", 100)
	if err != nil {
		return nil, err
	}

	var vv []Contact
	var v Contact
	for _, o := range out {
		_ = json.Unmarshal(o, &v)
		vv = append(vv, v)
	}

	return json.Marshal(&vv)
}

func handleDelete(ctx context.Context, idCsv string) ([]byte, error) {
	ids := strings.Split(idCsv, ",")
	if len(ids) == 0 {
		return nil, errors.New("ids are required")
	}

	client := s3.NewClient(ctx)
	var results = map[string]interface{}{
		"success": 0,
		"failure": 0,
	}
	for _, id := range ids {
		if err := client.Delete(ctx, bucket, "message/contact/read/"+id+".json"); err == nil {
			results[id] = true
			results["success"] = results["success"].(int) + 1
		} else {
			results[id] = false
			results["failure"] = results["failure"].(int) + 1
		}
	}
	return json.Marshal(results)
}

func handlePatch(ctx context.Context, id, read string) error {
	if len(id) == 0 {
		return errors.New("id is required")
	}
	unreadKey := "message/contact/unread/" + id + ".json"
	readKey := "message/contact/read/" + id + ".json"
	if read == "true" {
		return s3.NewClient(ctx).Move(ctx, bucket, unreadKey, readKey)
	}
	return s3.NewClient(ctx).Move(ctx, bucket, readKey, unreadKey)
}

func main() {
	lambda.Start(handler)
}
