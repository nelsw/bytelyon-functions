package model

import "testing"

func TestMakeURL(t *testing.T) {
	url := MakeURL("www.google.com/wat/")
	t.Log(url, url.Domain())
	err := url.Validate()
	t.Log(err)
	_, err = url.Document()
	t.Log(err)
}
