package auth

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(".env file not found")
	}
}

func TestNewAccessToken(t *testing.T) {

	data := ulid.Make().String()

	claims, err := Validate(NewToken(data))
	if err != nil {
		t.Fatal(err)
	}

	if data != claims.Data {
		t.Fatalf("got %s, want %s", claims.Data, data)
	}
}
