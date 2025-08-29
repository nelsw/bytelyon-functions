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
	id := ulid.Make()

	_ = New().Value(&FooBar{ID: id}).Save()

	var v FooBar
	if err := New(ctx).Value(&v).ID(id).Find(); err != nil {
		t.Error(err)
	}
	if v.ID != id {
		t.Errorf("want %s, got %s", id, v.ID)
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
