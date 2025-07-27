package logging

import (
	fxzerolog "github.com/efectn/fx-zerolog"
	server_config "github.com/greencoda/auth0-api-gateway/internal/config/server"
	"github.com/rs/zerolog"
	"go.uber.org/fx/fxevent"
)

func NewFXLogger(
	config *server_config.Config,
	logger zerolog.Logger,
) fxevent.Logger {
	if config.ReleaseStage != "local" {
		return fxevent.NopLogger
	}

	return &fxzerolog.ZeroLogger{Logger: logger}
}
