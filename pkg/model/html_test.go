package model

import (
	"bytelyon-functions/pkg/util/pretty"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
		FormatLevel: func(a any) string {
			var color, level string
			switch level = strings.ToUpper(a.(string)[:3]); level {
			case "TRA":
				color = "\033[36m" // Cyan
			case "DEB":
				color = "\033[35m" // Magenta
			case "INF":
				color = "\033[32m" // Green
			case "WAR":
				color = "\033[33m" // Yellow
			case "ERR":
				color = "\033[31m" // Red
			case "FAT", "PAN":
				color = "\033[41m\033[37m" // Red background, white text
			default:
				color = "\033[0m" // Reset
			}
			return color + level + "\033[0m" // Reset color after level
		},
	})
}

func TestMakeHTML(t *testing.T) {
	b, err := os.ReadFile("../../tests/fixtures/html/google_serp.html")
	if err != nil {
		panic(err)
	}

	html := MakeHTML(string(b))
	//for _, tag := range html.Tags("script") {
	//	fmt.Println(tag)
	//}

	//for idx, tag := range html.TagsWithClass(atom.Div, "top-pla-group-inner") {
	//fmt.Println(idx, tag)
	//}

	for _, tag := range html.TagsWithAttribute("div", "data-aavs") {
		fmt.Println(tag)
	}

}

func TestHTML_SerpProducts(t *testing.T) {
	b, err := os.ReadFile("../../tests/fixtures/html/google_serp.html")
	if err != nil {
		panic(err)
	}

	html := MakeHTML(string(b))

	pretty.Println(html.SerpSponsoredProducts())
}

func TestHTML_SerpOrganicResults(t *testing.T) {
	b, err := os.ReadFile("../../tests/fixtures/html/google_serp_ev_fire_blanket_for_sale_formatted.html")
	if err != nil {
		panic(err)
	}

	html := MakeHTML(string(b))
	for k, v := range html.SerpOrganicResults() {
		fmt.Printf("%v: {\n", k)
		for key, val := range v.(map[string]any) {
			if key == "img" {
				val = len(v.(map[string]any))
			}
			pretty.Println(map[string]any{key: val})
		}
		fmt.Printf("}\n")
	}
}
