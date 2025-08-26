package entity

import (
	"context"
	"fmt"
	"testing"

	"github.com/oklog/ulid/v2"
)

var ctx = context.Background()

type FooBar struct {
	ID ulid.ULID `json:"user_id"`
}

func TestSave(t *testing.T) {
	if err := New(ctx).Value(FooBar{ID: ulid.Make()}).Save(); err != nil {
		t.Error(err)
	}
}

func TestFind(t *testing.T) {
	v := FooBar{}
	if err := New(ctx).Value(&v).ID("01K30DA8S6HSF4F966V1BA67ZY").Find(); err != nil {
		t.Error(err)
	}
}

func TestPage(t *testing.T) {
	var v []FooBar
	if err := New(ctx).Path(&FooBar{}).Type(&v).Page(100); err != nil {
		t.Error(err)
	}
	fmt.Println(v)
}

func TestFoo(t *testing.T) {

}
