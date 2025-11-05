package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Service interface {
	Delete(string) error
	Find(string, any) error
	Get(string) ([]byte, error)
	Put(string, []byte) error
	Keys(string, string, int) ([]string, error)
}

type client struct {
	*s3.Client
	ctx    context.Context
	bucket string
}

func (c *client) Delete(k string) error {
	_, err := c.DeleteObject(c.ctx, &s3.DeleteObjectInput{
		Bucket: &c.bucket,
		Key:    key(k),
	})
	return err
}

func (c *client) Find(k string, a any) error {
	out, err := c.Get(k)
	if err != nil {
		return err
	}
	return json.Unmarshal(out, &a)
}

func (c *client) Get(k string) (b []byte, err error) {
	var out *s3.GetObjectOutput
	out, err = c.GetObject(c.ctx, &s3.GetObjectInput{
		Bucket: &c.bucket,
		Key:    key(k),
	})
	if err == nil {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(out.Body)
		b, err = io.ReadAll(out.Body)
	}
	return
}

func (c *client) Put(k string, data []byte) error {
	_, err := c.PutObject(c.ctx, &s3.PutObjectInput{
		Bucket: &c.bucket,
		Key:    key(k),
		Body:   bytes.NewReader(data),
	})
	return err
}

func (c *client) Keys(prefix, after string, size int) (keys []string, err error) {
	maxKeys := int32(size)
	if maxKeys == 0 {
		maxKeys = 10
	}
	input := s3.ListObjectsV2Input{
		Bucket:  &c.bucket,
		Prefix:  &prefix,
		MaxKeys: &maxKeys,
	}
	if after != "" {
		input.StartAfter = &after
	}
	var out *s3.ListObjectsV2Output
	if out, err = c.ListObjectsV2(c.ctx, &input); err == nil {
		for _, o := range out.Contents {
			keys = append(keys, *o.Key)
		}
	}
	return
}

func key(s string) *string {
	if strings.HasPrefix(s, "/") {
		s = s[1:]
	}
	if !strings.HasSuffix(s, "/_.json") {
		s += "/_.json"
	}
	return &s
}

// New returns a new S3 service with the provided context.
func New() Service {
	ctx := context.Background()
	cfg, _ := config.LoadDefaultConfig(ctx)
	mode := os.Getenv("APP_MODE")
	if mode == "" {
		mode = "test"
	}
	return &client{
		s3.NewFromConfig(cfg),
		ctx,
		"bytelyon-db-" + mode,
	}
}
