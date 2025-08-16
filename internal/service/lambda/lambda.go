package lambda

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

var exc *types.ResourceConflictException

type Service interface {
	Create(ctx context.Context, name, role string, zipFile []byte, vars map[string]string)
	Update(ctx context.Context, name string, zipFile []byte, vars map[string]string)
	Delete(ctx context.Context, name string)
	Publish(ctx context.Context, name string)
}

type Client struct {
	*lambda.Client
}

func NewClient(ctx context.Context) Service {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return &Client{lambda.NewFromConfig(cfg)}
}

func (c *Client) Create(ctx context.Context, name, role string, zipFile []byte, vars map[string]string) {

	if _, err := c.CreateFunction(ctx, &lambda.CreateFunctionInput{
		Code:          &types.FunctionCode{ZipFile: zipFile},
		FunctionName:  &name,
		Role:          &role,
		Handler:       aws.String("bootstrap"),
		Publish:       true,
		Runtime:       types.RuntimeProvidedal2,
		Architectures: []types.Architecture{types.ArchitectureArm64},
		Timeout:       aws.Int32(int32(90)),
		Environment:   &types.Environment{Variables: vars},
	}); err == nil {
		_, err = lambda.
			NewFunctionActiveV2Waiter(c).
			WaitForOutput(ctx, &lambda.GetFunctionInput{FunctionName: &name}, time.Second*15)
	} else {
		log.Panicf("Create function failed, %v", err)
	}
}

func (c *Client) Update(ctx context.Context, name string, zipFile []byte, vars map[string]string) {

	if vars != nil {
		if _, err := c.UpdateFunctionConfiguration(ctx, &lambda.UpdateFunctionConfigurationInput{
			FunctionName: &name,
			Environment:  &types.Environment{Variables: vars},
		}); err != nil {
			log.Panicf("Update function failed, %v", err)
		}
		time.Sleep(30 * time.Second)
	}

	if _, err := c.UpdateFunctionCode(ctx, &lambda.UpdateFunctionCodeInput{
		FunctionName: &name,
		ZipFile:      zipFile,
		Publish:      true,
	}); err != nil {
		log.Panicf("Update function failed, %v", err)
	}
}

func (c *Client) Publish(ctx context.Context, name string) {
	if _, err := c.CreateFunctionUrlConfig(ctx, &lambda.CreateFunctionUrlConfigInput{
		AuthType:     types.FunctionUrlAuthTypeNone,
		FunctionName: &name,
	}); err != nil {
		log.Panicf("CreateFunctionUrlConfig failed, %v", err)
	}

	if _, err := c.AddPermission(ctx, &lambda.AddPermissionInput{
		Action:              aws.String("lambda:InvokeFunctionUrl"),
		FunctionName:        &name,
		Principal:           aws.String("*"),
		StatementId:         aws.String("FunctionURLAllowPublicAccess"),
		FunctionUrlAuthType: types.FunctionUrlAuthTypeNone,
	}); err != nil {
		log.Panicf("AddPermission failed, %v", err)
	}
}

func (c *Client) Delete(ctx context.Context, name string) {
	if _, err := c.GetFunction(ctx, &lambda.GetFunctionInput{FunctionName: &name}); err == nil {
		if _, err = c.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: &name}); err != nil {
			log.Panicf("Delete function failed, %v", err)
		}
	}
}
