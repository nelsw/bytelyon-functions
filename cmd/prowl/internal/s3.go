package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3 Service

func init() {
	InitLogger()
	cfg, _ := config.LoadDefaultConfig(context.Background())
	S3 = &client{
		s3.NewFromConfig(cfg),
		context.Background(),
		"bytelyon-db-prod",
	}
}

type Service interface {
	Get(string) ([]byte, error)
	Find(any) error
	Put(any, ...any) error
	Keys(string) ([]string, error)
}

type client struct {
	*s3.Client
	context.Context
	bucket string
}

func (c *client) Get(k string) ([]byte, error) {
	out, err := c.GetObject(c.Context, &s3.GetObjectInput{
		Bucket: &c.bucket,
		Key:    &k,
	})
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()
	return io.ReadAll(out.Body)
}

func (c *client) Find(a any) error {
	b, err := S3.Get(fmt.Sprint(a) + "/_.json")

	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, a); err != nil {
		return err
	}
	return nil
}

func (c *client) Put(a any, aa ...any) error {
	name := "/_.json"
	if len(aa) > 1 {
		name = aa[1].(string)
	}
	var b []byte
	if len(aa) > 0 {
		b = Bytes(aa[0])
	} else {
		b = Bytes(a)
	}
	_, err := c.PutObject(c.Context, &s3.PutObjectInput{
		Bucket: &c.bucket,
		Key:    Ptr(fmt.Sprint(a) + name),
		Body:   bytes.NewReader(b),
	})
	return err
}

func (c *client) Keys(k string) ([]string, error) {
	out, err := c.ListObjectsV2(c.Context, &s3.ListObjectsV2Input{
		Bucket:  &c.bucket,
		Prefix:  &k,
		MaxKeys: aws.Int32(1000),
	})
	if err != nil {
		return nil, err
	}
	var keys []string
	for _, o := range out.Contents {
		keys = append(keys, *o.Key)
	}
	return keys, nil
}
