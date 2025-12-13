package api

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
		FieldsOrder: []string{
			"user",
			"user_id",
			"search",
			"search_id",
			"job",
			"job_id",
			"result",
			"result_id",
			"page",
			"page_id",
			"id",
			"plunder",
			"news",
			"target",
			"follow",
		},
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
