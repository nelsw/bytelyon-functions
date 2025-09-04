package db

import (
	"bytelyon-functions/pkg/service/s3"
	"context"
	"encoding/json"
	"os"
	"strings"
)

const bucket = "bytelyon"

var ctx context.Context
var client s3.Service
var db string

func init() {
	ctx = context.Background()
	client = s3.NewClient(ctx)
	mode, ok := os.LookupEnv("APP_MODE")
	if mode == "" || !ok {
		mode = "local"
	}
	db = mode + "/db"
}

func key(s string) string {
	if !strings.HasSuffix(s, ".json") {
		s += ".json"
	}
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}
	return db + s
}

func Save(k string, a any) error {
	b, _ := json.Marshal(&a)
	return client.Put(ctx, bucket, key(k), b)
}

func FindOne(k string, a any) error {
	b, err := GetOne(k)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &a)
}

func GetOne(k string) ([]byte, error) {
	return client.Get(ctx, bucket, key(k))
}
