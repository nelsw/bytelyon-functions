package util

import (
	"encoding/json"
	"os"
)

func AppMode() string {
	mode := os.Getenv("APP_MODE")
	if mode == "" {
		mode = "local"
	}
	return mode
}

func First(a ...any) any {
	if a != nil && len(a) == 1 || a[0] != nil {
		return a[0]
	}
	return nil
}

func IsJSON(s string) bool {
	if s == "" {
		return false
	}
	var raw json.RawMessage
	return json.Unmarshal([]byte(s), &raw) == nil
}

func MustMarshal(a any) []byte {
	b, err := json.Marshal(&a)
	if err != nil {
		panic(err)
	}
	return b
}

func MustUnmarshal(b []byte, a any) {
	if err := json.Unmarshal(b, &a); err != nil {
		panic(err)
	}
}
