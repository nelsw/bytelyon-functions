package helper

import (
	"encoding/json"
	"fmt"
	"regexp"
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

func SplitStringByCase(s string) string {
	reg := regexp.MustCompile("([A-Z]+)")
	return reg.ReplaceAllString(s, `/$1`)
}
