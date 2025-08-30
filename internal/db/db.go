package db

import (
	"bytelyon-functions/pkg/service/s3"
	"context"
	"encoding/json"
)

const bucket = "bytelyon"
const object = ".json"

var ctx context.Context
var client s3.Service

func init() {
	ctx = context.Background()
	client = s3.NewClient(ctx)
}

func Save(k string, a any) error {
	b, _ := json.Marshal(&a)
	return client.Put(ctx, bucket, k+object, b)
}

func Get(k string, a any) error {
	b, err := client.Get(ctx, bucket, k+object)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &a)
}
