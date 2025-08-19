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
	out, err := NewToken(data)
	if err != nil {
		t.Fatal(err)
	}

	var claims *Claims
	claims, err = Validate(out)
	if err != nil {
		t.Fatal(err)
	}
	if claims.Data != data {
		t.Fatalf("got %s, want %s", claims.Data, data)
	}
}
