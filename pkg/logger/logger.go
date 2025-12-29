package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = *New()
}

type builder struct {
	level  zerolog.Level
	fields []string
	caller bool
}

func New(a ...any) *zerolog.Logger {
	b := builder{
		level: zerolog.DebugLevel,
	}
	for _, v := range a {
		switch v.(type) {
		case zerolog.Level:
			b.level = v.(zerolog.Level)
		case []string:
			b.fields = v.([]string)
		case bool:
			b.caller = v.(bool)
		}
	}
	return b.build()
}

func (b builder) build() *zerolog.Logger {
	l := log.Output(zerolog.ConsoleWriter{
		Out:         os.Stdout,
		FieldsOrder: b.fields,
		FormatLevel: func(a any) string {
			switch l := strings.ToUpper(a.(string)[:3]); l {
			case "TRA":
				return Cyan + l + Default
			case "DEB":
				return Purple + l + Default
			case "INF":
				return Green + l + Default
			case "WAR":
				return Yellow + l + Default
			case "ERR":
				return Red + l + Default
			case "FAT", "PAN":
				return RedBackground + White + l + Default
			default:
				return Default + l + Default
			}
		},
	}).Level(b.level)
	if b.caller {
		l = l.With().Caller().Logger()
	}
	return &l
}
