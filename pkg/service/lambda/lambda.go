package lambda

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type Service interface {
	InvokeEvent(context.Context, string, []byte) ([]byte, error)
	InvokeRequest(context.Context, string, []byte) ([]byte, error)
}

type Client struct {
	*lambda.Client
}

func New(ctx context.Context) Service {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return &Client{lambda.NewFromConfig(cfg)}
}

func (c *Client) InvokeEvent(ctx context.Context, name string, payload []byte) ([]byte, error) {
	return c.invoke(ctx, types.InvocationTypeEvent, name, payload)
}

func (c *Client) InvokeRequest(ctx context.Context, name string, payload []byte) ([]byte, error) {
	return c.invoke(ctx, types.InvocationTypeRequestResponse, name, payload)
}

func (c *Client) invoke(ctx context.Context, typ types.InvocationType, name string, payload []byte) ([]byte, error) {

	output, err := c.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   &name,
		InvocationType: typ,
		Payload:        payload,
	})

	if err != nil {
		return nil, err
	}

	return output.Payload, nil
}
