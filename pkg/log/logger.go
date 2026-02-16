package log

import (
	stdlog "log"
	"os"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/subosito/gotenv"
)

var Logger zerolog.Logger

func init() {
	err := gotenv.Load()
	if err != nil {
		stdlog.Printf("MISSING ENVS: %v", err)
	}

	err = loadSpecs()
	if err != nil {
		stdlog.Fatalf("logger env parse error: %v", err)
	}

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"

	initLevel()
	initFormat()

	stdlog.SetFlags(0)
	stdlog.SetOutput(Logger)

	Logger.Info().Msgf("logger initialized (%s level)", zerolog.GlobalLevel())
}

func initFormat() {
	if specs.LogFormat == "json" {
		Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		consoleWriter := zerolog.ConsoleWriter{
			Out:         os.Stderr,
			TimeFormat:  specs.TimeFormat,
			FieldsOrder: []string{"source", "component"},
		}
		Logger = zlog.Output(consoleWriter)
	}
}

func initLevel() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if specs.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func For(source string) *zerolog.Logger {
	logger := Logger.With().Str("source", source).Logger()
	return &logger
}
