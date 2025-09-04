package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Service interface {
	Delete(ctx context.Context, bucket, key string) error
	Get(ctx context.Context, bucket, key string) ([]byte, error)
	Put(ctx context.Context, bucket, key string, data []byte) error
	List(ctx context.Context, bucket, prefix string, maxKeys int32) ([][]byte, error)
	Move(ctx context.Context, bucket, oldKey, newKey string) error
	Keys(ctx context.Context, size int32, bucket, prefix string) ([]string, error)
	KeysAfter(ctx context.Context, size int32, bucket, prefix, after string) ([]string, error)
}

type Client struct {
	*s3.Client
}

func (c *Client) Move(ctx context.Context, bucket, oldKey, newKey string) error {

	b, err := c.Get(ctx, bucket, oldKey)
	if err != nil {
		return err
	}

	if err = c.Put(ctx, bucket, newKey, b); err != nil {
		return err
	}

	return c.Delete(ctx, bucket, oldKey)
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
	fmt.Println(out.Metadata)
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

func (c *Client) List(ctx context.Context, bucket, prefix string, maxKeys int32) ([][]byte, error) {
	output, err := c.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:  &bucket,
		Prefix:  &prefix,
		MaxKeys: &maxKeys,
	})
	if err != nil {
		return nil, err
	}

	var results [][]byte
	for _, out := range output.Contents {
		o, e := c.Get(ctx, bucket, *out.Key)
		if e != nil {
			log.Println("Get Error:", e)
		} else {
			results = append(results, o)
		}
	}

	return results, nil
}

func (c *Client) Keys(ctx context.Context, size int32, bucket, prefix string) ([]string, error) {
	return c.KeysAfter(ctx, size, bucket, prefix, "")
}

func (c *Client) KeysAfter(ctx context.Context, size int32, bucket, prefix, after string) ([]string, error) {
	input := s3.ListObjectsV2Input{
		Bucket:  &bucket,
		Prefix:  &prefix,
		MaxKeys: &size,
	}
	if after != "" {
		input.StartAfter = &after
	}
	output, err := c.ListObjectsV2(ctx, &input)
	if err != nil {
		return nil, err
	}
	var keys []string
	for _, out := range output.Contents {
		keys = append(keys, *out.Key)
	}
	return keys, nil
}

func NewClient(ctx context.Context) Service {
	if cfg, err := config.LoadDefaultConfig(ctx); err != nil {
		panic(err)
	} else {
		return &Client{s3.NewFromConfig(cfg)}
	}
}
