package db

import (
	"bytelyon-functions/pkg/util"
	"bytes"
	"context"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

type S3 interface {
	Delete(string) error
	Get(string) ([]byte, error)
	Put(string, []byte) error
	Keys(...string) ([]string, error)
	URL(string, int64) (string, error)
}

type s3 struct {
	Bucket *string
	*_s3.Client
	*_s3.PresignClient
	context.Context
}

func (s3 *s3) Delete(k string) error {

	_, err := s3.DeleteObject(s3.Context, &_s3.DeleteObjectInput{
		Bucket: s3.Bucket,
		Key:    &k,
	})

	log.Trace().
		Err(err).
		Str("key", k).
		Msg("Delete")

	return err
}

func (s3 *s3) Get(k string) ([]byte, error) {

	out, err := s3.GetObject(s3.Context, &_s3.GetObjectInput{
		Bucket: s3.Bucket,
		Key:    util.Ptr(k),
	})
	if err != nil {
		log.Warn().Err(err).Msg("Get - Failed to read body")
		return nil, err
	}
	var body []byte
	defer out.Body.Close()
	if body, err = io.ReadAll(out.Body); err != nil {
		log.Warn().Err(err).Msg("Get - Failed to read body")
		return nil, err
	}

	log.Trace().
		Err(err).
		Str("key", k).
		Bytes("body", body).
		Msg("Get")

	return body, err
}

func (s3 *s3) Put(k string, b []byte) (err error) {

	_, err = s3.PutObject(s3.Context, &_s3.PutObjectInput{
		Bucket: s3.Bucket,
		Key:    &k,
		Body:   bytes.NewReader(b),
	})

	log.Trace().
		Err(err).
		Str("key", k).
		Bytes("body", b).
		Msg("Put")

	return
}

func (s3 *s3) Keys(s ...string) ([]string, error) {

	var keys []string
	var err error

	var fn func(string, string)
	fn = func(prefix, after string) {
		var out *_s3.ListObjectsV2Output
		out, err = s3.ListObjectsV2(s3.Context, &_s3.ListObjectsV2Input{
			Bucket:     s3.Bucket,
			Prefix:     &prefix,
			MaxKeys:    util.Ptr(int32(1000)),
			StartAfter: &after,
		})

		if err == nil {
			for _, obj := range out.Contents {
				keys = append(keys, *obj.Key)
			}
		}

		if len(out.Contents) == 1000 {
			fn(prefix, *out.Contents[len(out.Contents)-1].Key)
		}
	}

	if len(s) == 1 {
		s = append(s, "")
	}

	fn(s[0], s[1])

	log.Debug().
		Err(err).
		Str("prefix", s[0]).
		Str("after", s[1]).
		Int("keys", len(keys)).
		Msg("Keys")

	return keys, err
}

func (s3 *s3) URL(k string, i int64) (string, error) {

	out, err := s3.PresignGetObject(s3.Context, &_s3.GetObjectInput{
		Bucket: s3.Bucket,
		Key:    &k,
	}, _s3.WithPresignExpires(time.Duration(i)*time.Minute))

	var url string
	if out != nil {
		url = out.URL
	}

	log.Trace().
		Err(err).
		Str("key", k).
		Int64("exp", i).
		Str("url", url).
		Msg("URL")

	return url, err
}

// NewS3 returns a new S3 client with a background context.
// An optional variadic set of Config values can be provided as
// input that will be prepended to the configs slice.
func NewS3(optFns ...func(*config.LoadOptions) error) S3 {
	return NewS3WithContext(context.Background(), optFns...)
}

// NewS3WithContext returns a new S3 client with the provided context.
// An optional variadic set of Config values can be provided as
// input that will be prepended to the configs slice.
func NewS3WithContext(ctx context.Context, optFns ...func(*config.LoadOptions) error) S3 {
	cfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		panic(err)
	}
	b := os.Getenv("S3_BUCKET")
	if b == "" {
		panic("S3_BUCKET environment variable must be set")
	}
	c := _s3.NewFromConfig(cfg)
	return &s3{
		&b,
		c,
		_s3.NewPresignClient(c),
		ctx,
	}
}
