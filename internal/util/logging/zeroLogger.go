package logging

import (
	"os"
	"strings"

	server_config "github.com/greencoda/auth0-api-gateway/internal/config/server"
	"github.com/rs/zerolog"
)

func NewZeroLogger(config *server_config.Config) zerolog.Logger {
	logger := setLogLevel(
		zerolog.New(os.Stdout).
			With().
			Timestamp().
			Logger(),
		config.LogLevel,
	)

	if config.ReleaseStage != "local" {
		return logger
	}

	return logger.Output(zerolog.ConsoleWriter{
		TimeFormat: zerolog.TimeFieldFormat,
		Out:        os.Stderr,
	})
}

func setLogLevel(logger zerolog.Logger, logLevel string) zerolog.Logger {
	switch strings.ToLower(logLevel) {
	case "debug":
		return logger.Level(zerolog.DebugLevel)
	case "info":
		return logger.Level(zerolog.InfoLevel)
	case "warn":
		return logger.Level(zerolog.WarnLevel)
	case "error":
		return logger.Level(zerolog.ErrorLevel)
	case "fatal":
		return logger.Level(zerolog.FatalLevel)
	case "panic":
		return logger.Level(zerolog.PanicLevel)
	}

	return logger
}
