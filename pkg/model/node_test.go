package model

import (
	"encoding/json"
	"testing"
)

func TestNewNode(t *testing.T) {
	m := map[string]string{"a": "b"}
	n := NewNode("id", "label", m)
	b, _ := json.MarshalIndent(n, "", "\t")
	t.Log(string(b))
}
