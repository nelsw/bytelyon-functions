package main

import (
	"bytelyon-functions/test"
	"encoding/base64"
	"testing"
)

func TestLogin(t *testing.T) {
	e := "kowalski7012@gmail.com"
	p := "Farts1234!"
	s := base64.StdEncoding.EncodeToString([]byte(e + ":" + p))

	test.New(t).
		Header("authorization", s).Post(nil).
		Handle(handler).
		OK().JSON(map[string]any{})
}
