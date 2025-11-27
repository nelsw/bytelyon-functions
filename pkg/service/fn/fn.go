package fn

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type Service interface {
	Request(string, any) ([]byte, error)
	Event(string, any) error
}

type Client struct {
	*lambda.Client
	context.Context
}

func (c *Client) Request(name string, a any) ([]byte, error) {
	return c.invoke(name, a, types.InvocationTypeEvent)
}

func (c *Client) Event(name string, a any) error {
	_, err := c.invoke(name, a, types.InvocationTypeEvent)
	return err
}

func (c *Client) invoke(name string, a any, t types.InvocationType) ([]byte, error) {

	b, err := json.Marshal(&a)
	if err != nil {
		return nil, err
	}

	var out *lambda.InvokeOutput
	out, err = c.Invoke(c.Context, &lambda.InvokeInput{
		FunctionName:   &name,
		InvocationType: t,
		Payload:        b,
	})

	if err != nil {
		return nil, err
	}

	return out.Payload, nil
}

func New(ctx context.Context) Service {
	ctx = context.Background()
	cfg, _ := config.LoadDefaultConfig(ctx)
	return &Client{lambda.NewFromConfig(cfg), ctx}
}
