package pretty

import (
	"encoding/json"
	"fmt"
)

func Of(a any) string {
	b, _ := json.MarshalIndent(a, "", "  ")
	return string(b)
}

func Println(a any) {
	fmt.Println(Of(a))
}
