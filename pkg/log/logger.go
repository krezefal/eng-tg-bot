package log

import "github.com/rs/zerolog"

var Logger zerolog.Logger

func For(source string) *zerolog.Logger {
	logger := Logger.With().Str("source", source).Logger()
	return &logger
}
