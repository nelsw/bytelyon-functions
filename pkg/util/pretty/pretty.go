package pretty

import (
	"encoding/json"
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var p = message.NewPrinter(language.English)

func Of(a any) string {
	b, _ := json.MarshalIndent(a, "", "  ")
	return string(b)
}

func Println(a any) {
	fmt.Println(Of(a))
}

func Number(i int) string {
	return p.Sprintf("%d", i)
}
