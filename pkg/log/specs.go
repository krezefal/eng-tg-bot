package logger

import (
	"os"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type Specs struct {
	Debug      bool   `envconfig:"DEBUG" default:"false"`
	LogFormat  string `envconfig:"LOG_FORMAT" default:"console" example:"console, json"`
	TimeFormat string `envconfig:"TIME_FORMAT" default:"2006-01-02 15:04:05"`
}

var specs Specs

func initFormat() {
	if specs.LogFormat == "json" {
		Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: specs.TimeFormat})
	}
}
