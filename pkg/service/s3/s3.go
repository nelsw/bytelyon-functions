package s3

import (
	"bytes"
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Service interface {
	Delete(ctx context.Context, bucket, key string) error
	Get(ctx context.Context, bucket, key string) ([]byte, error)
	Put(ctx context.Context, bucket, key string, data []byte) error
	List(ctx context.Context, bucket string, prefix *string) ([]string, error)
}

type Client struct {
	*s3.Client
}

func NewClient(ctx context.Context) Service {
	if cfg, err := config.LoadDefaultConfig(ctx); err != nil {
		panic(err)
	} else {
		return &Client{s3.NewFromConfig(cfg)}
	}
}

func (c *Client) Delete(ctx context.Context, bucket, key string) error {
	_, err := c.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	return err
}

func (c *Client) Get(ctx context.Context, bucket, key string) ([]byte, error) {
	out, err := c.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Fatal(err)
		}
	}(out.Body)

	return io.ReadAll(out.Body)
}

func (c *Client) Put(ctx context.Context, bucket, key string, data []byte) error {
	_, err := c.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader(data),
	})
	return err
}

func (c *Client) List(ctx context.Context, bucket string, prefix *string) ([]string, error) {
	output, err := c.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: prefix,
	})
	if err != nil {
		return nil, err
	}

	var keys []string
	for _, object := range output.Contents {
		keys = append(keys, *object.Key)
	}

	return keys, nil
}
