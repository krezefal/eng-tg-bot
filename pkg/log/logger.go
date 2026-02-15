package logger

import (
	"log"

	"github.com/rs/zerolog"
	"github.com/subosito/gotenv"
)

var Logger zerolog.Logger

func init() {
	if err := gotenv.Load(); err != nil {
		For("logger").Warn().Err(err).Msg("missing envs")
	}

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"

	initLevel()
	initFormat()

	log.SetFlags(0)
	log.SetOutput(Logger)

	Logger.Info().Msgf("logger initialized (%s level)", zerolog.GlobalLevel())
}

func For(source string) *zerolog.Logger {
	logger := Logger.With().Str("source", source).Logger()
	return &logger
}
