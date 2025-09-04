package ses

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
)

func TestClient_VerifyEmail(t *testing.T) {
	ctx := context.Background()
	to := "kowalski7012@gmail.com"
	url := gofakeit.URL()

	if err := New(ctx).VerifyEmail(ctx, to, url); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ResetPassword(t *testing.T) {
	ctx := context.Background()
	to := "kowalski7012@gmail.com"
	url := gofakeit.URL()

	if err := New(ctx).ResetPassword(ctx, to, url); err != nil {
		t.Fatal(err)
	}
}
