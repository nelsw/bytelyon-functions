package util

import (
	"encoding/json"
)

func IsJSON(s string) bool {
	if s == "" {
		return false
	}
	var raw json.RawMessage
	return json.Unmarshal([]byte(s), &raw) == nil
}

func First(a ...any) any {
	if a != nil && len(a) == 1 || a[0] != nil {
		return a[0]
	}
	return nil
}

func MustMarshal(a any) []byte {
	b, err := json.Marshal(&a)
	if err != nil {
		panic(err)
	}
	return b
}
