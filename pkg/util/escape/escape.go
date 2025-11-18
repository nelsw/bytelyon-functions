package escape

import "strings"

var hexadecimals = map[string]string{
	"\\x22": "\"",
	"\\x26": "&",
	"\\x27": "'",
	"\\x3c": "<",
	"\\x3d": "=",
	"\\x3e": ">",
	"\\xb7": ".",
}

var unicodeChars = map[string]string{
	"\u003c": "<",
	"\u003e": ">",
}

var htmlEntities = map[string]string{
	"&#37;": "%",
	"&amp;": "&",
}

func Hexadecimal(s string) string {
	for k, v := range hexadecimals {
		if strings.Contains(s, k) {
			s = strings.ReplaceAll(s, k, v)
		}
	}
	return s
}

func Unicode(s string) string {
	for k, v := range unicodeChars {
		if strings.Contains(s, k) {
			s = strings.ReplaceAll(s, k, v)
		}
	}
	return s
}

func HtmlEntity(s string) string {
	for k, v := range htmlEntities {
		if strings.Contains(s, k) {
			s = strings.ReplaceAll(s, k, v)
		}
	}
	return s
}
