package helper

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type UserEmail struct{}

func TestSplitStringByCase(t *testing.T) {
	var ue UserEmail
	var name string
	if typ := reflect.TypeOf(ue); typ.Kind() == reflect.Ptr {
		name = typ.Elem().Name()
	} else {
		name = typ.Name()
	}
	ss := SplitStringByCase(name)
	ss = strings.ToLower(ss)
	fmt.Println(ss)
}
