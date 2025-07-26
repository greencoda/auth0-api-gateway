package logging

import (
	fxzerolog "github.com/efectn/fx-zerolog"
	server_config "github.com/greencoda/auth0-api-gateway/internal/config/server"
	"github.com/rs/zerolog"
	"go.uber.org/fx/fxevent"
)

// NewFXLogger returns a new fxevent.Logger for the FX lifecycle events
func NewFXLogger(
	config *server_config.Config,
	logger zerolog.Logger,
) fxevent.Logger {
	if config.ReleaseStage == "prod" {
		return fxevent.NopLogger
	}

	return &fxzerolog.ZeroLogger{Logger: logger}
}
