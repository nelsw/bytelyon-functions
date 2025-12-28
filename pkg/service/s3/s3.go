package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

type Service interface {
	Delete(string) error
	Find(string, any) error
	Get(string) ([]byte, error)
	GetPresigned(string) (string, error)
	Put(string, []byte) error
	Keys(string, string, int) ([]string, error)
}

type client struct {
	*s3.Client
	*s3.PresignClient
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

func (c *client) Put(k string, data []byte) (err error) {
	_, err = c.PutObject(c.ctx, &s3.PutObjectInput{
		Bucket: &c.bucket,
		Key:    key(k),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		log.Warn().Err(err).Str("key", k).Bytes("body", data).Msg("s3 - Failed Put")
	}
	return
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

func (c *client) GetPresigned(k string) (string, error) {
	out, err := c.PresignGetObject(c.ctx, &s3.GetObjectInput{
		Bucket: &c.bucket,
		Key:    &k,
	}, s3.WithPresignExpires(time.Duration(30*int64(time.Minute))))
	if err != nil {
		return "", err
	}
	return out.URL, nil
}

func key(s string) *string {
	if strings.HasPrefix(s, "/") {
		s = s[1:]
	}
	if strings.HasSuffix(s, ".html") ||
		strings.HasSuffix(s, ".png") ||
		strings.HasSuffix(s, ".json") {
		return &s
	}
	if !strings.HasSuffix(s, "/_.json") {
		s += "/_.json"
	}
	return &s
}

// New returns a new S3 service with the provided context.
func New() Service {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}
	mode := os.Getenv("APP_MODE")
	if mode == "" {
		mode = "test"
	}
	c := s3.NewFromConfig(cfg)
	return &client{
		c,
		s3.NewPresignClient(c),
		context.Background(),
		"bytelyon-db-prod",
	}
}
