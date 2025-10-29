package lambda

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type Service interface {
	Request(context.Context, string, any) ([]byte, error)
}

type Client struct {
	*lambda.Client
}

func (c *Client) Request(ctx context.Context, name string, a any) (out []byte, err error) {

	var in []byte
	if in, err = json.Marshal(&a); err != nil {
		return nil, err
	}

	var output *lambda.InvokeOutput
	output, err = c.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   &name,
		InvocationType: types.InvocationTypeRequestResponse,
		Payload:        in,
	})

	if err == nil {
		out = output.Payload
	}

	return
}

// New returns a new Lambda client with the provided context.
func New(ctx context.Context) Service {
	cfg, _ := config.LoadDefaultConfig(ctx)
	return &Client{lambda.NewFromConfig(cfg)}
}
