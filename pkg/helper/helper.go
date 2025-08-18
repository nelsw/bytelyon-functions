package helper

import (
	"encoding/json"
	"fmt"
)

func Marshal(v interface{}) string {
	b, _ := json.Marshal(&v)
	return string(b)
}

func Pretty(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "\t")
	return string(b)
}

func PrettyPrintln(v interface{}) {
	fmt.Println(Pretty(v))
}

func PrintlnJson(v interface{}) {
	fmt.Println(Marshal(v))
}
