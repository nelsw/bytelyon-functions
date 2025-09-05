package contact

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/s3"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/oklog/ulid/v2"
)

type Contact struct {
	ID    ulid.ULID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Value string    `json:"message"`
}

func (c Contact) Validate() (err error) {
	if c.Name == "" {
		err = errors.Join(err, errors.New("name is required"))
	}
	if c.Email == "" {
		err = errors.Join(err, errors.New("email is required"))
	}
	if c.Value == "" {
		err = errors.Join(err, errors.New("message is required"))
	}
	return
}

func (c Contact) Key() string {
	return fmt.Sprintf("message/contact/unread/%s", c.ID)
}

func Handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	app.LogURLRequest(req)

	if app.IsOptions(req) {
		return app.OK()
	}

	if app.IsPost(req) {
		var c Contact
		if err := json.Unmarshal([]byte(req.Body), &c); err == nil {
			err = c.Validate()
		}
		c.ID = app.NewUlid()
		return app.Err(s3.New(ctx).Put(c.Key(), app.MustMarshal(c)))
	}

	return app.NotImplemented(req)
}
