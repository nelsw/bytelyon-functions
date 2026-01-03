package db

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:         os.Stdout,
		FieldsOrder: []string{"key", "prefix", "after", "size", "body", "keys"},
	})
}

func testKey(id ...ulid.ULID) string {
	var ID = ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0")
	if len(id) > 0 {
		ID = id[0]
	}
	return "users/" + ID.String() + "/_.json"
}

func testBody(id ...ulid.ULID) []byte {
	var ID = ulid.MustParse("01K48PC0BK13BWV2CGWFP8QQH0")
	if len(id) > 0 {
		ID = id[0]
	}
	return []byte(`{"id":"` + ID.String() + `"}`)
}

func TestClient_Put(t *testing.T) {
	t.Setenv("S3_BUCKET", "bytelyon-db")
	assert.NoError(t, NewS3().Put(testKey(), testBody()))
}

func TestClient_Get(t *testing.T) {
	t.Setenv("S3_BUCKET", "bytelyon-db")
	out, err := NewS3().Get(testKey())
	assert.NoError(t, err)
	assert.Equal(t, testBody(), string(out))
}

func TestClient_Delete(t *testing.T) {
	t.Setenv("S3_BUCKET", "bytelyon-db")
	assert.NoError(t, NewS3().Delete(testKey()))
}

func TestClient_Keys(t *testing.T) {

	t.Setenv("S3_BUCKET", "bytelyon-db")

	var ids []ulid.ULID
	for i := 0; i < 10; i++ {
		ids = append(ids, ulid.Make())
		assert.NoError(t, NewS3().Put(testKey(ids[i]), testBody(ids[i])))
	}

	keys, err := NewS3().Keys("users/", testKey(ids[5]), "2")

	assert.NoError(t, err)
	assert.Len(t, keys, 2)
	assert.Equal(t, testKey(ids[6]), keys[0])
	assert.Equal(t, testKey(ids[7]), keys[1])

	for _, id := range ids {
		_ = NewS3().Delete(testKey(id))
	}
}

func TestClient_URL(t *testing.T) {

	t.Setenv("S3_BUCKET", "bytelyon-db")

	assert.NoError(t, NewS3().Put(testKey(), testBody()))

	url, err := NewS3().URL(testKey(), 5)
	assert.NoError(t, err)

	var out *http.Response
	out, err = http.Get(url)
	assert.NoError(t, err)

	defer out.Body.Close()
	var b []byte
	b, err = io.ReadAll(out.Body)

	assert.NoError(t, err)
	assert.Equal(t, testBody(), string(b))

	_ = NewS3().Delete(testKey())
}
