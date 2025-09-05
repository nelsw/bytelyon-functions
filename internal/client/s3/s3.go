package s3

import (
	"bytelyon-functions/internal/app"
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client interface {
	Delete(string) error
	Get(string) ([]byte, error)
	Put(string, []byte) error
	Move(string, string) error
	Keys(string, string, int) ([]string, error)
}

type client struct {
	*s3.Client
	ctx context.Context
}

func (c *client) Move(oldKey, newKey string) error {

	b, err := c.Get(oldKey)
	if err != nil {
		return err
	}

	if err = c.Put(newKey, b); err != nil {
		return err
	}

	return c.Delete(oldKey)
}

func (c *client) Delete(k string) error {
	_, err := c.DeleteObject(c.ctx, &s3.DeleteObjectInput{
		Bucket: app.Bucket(),
		Key:    key(k),
	})
	return err
}

func (c *client) Get(k string) (b []byte, err error) {
	var out *s3.GetObjectOutput
	out, err = c.GetObject(c.ctx, &s3.GetObjectInput{
		Bucket: app.Bucket(),
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
		Bucket: app.Bucket(),
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
		Bucket:  app.Bucket(),
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

// New returns a new S3 client with the default context.
func New() Client {
	return NewWithContext(context.Background())
}

// NewWithContext returns a new S3 client with the provided context.
func NewWithContext(ctx context.Context) Client {
	if cfg, err := config.LoadDefaultConfig(ctx); err != nil {
		panic(err)
	} else {
		return &client{s3.NewFromConfig(cfg), ctx}
	}
}
