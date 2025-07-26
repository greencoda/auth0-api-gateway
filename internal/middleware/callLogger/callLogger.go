package callLogger

import (
	"net/http"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

// ICallLogger interface defines the method to get the call logging middleware handler.
type ICallLogger interface {
	Handler(h http.Handler) http.Handler
}

// CallLogger implements the ICallLog interface and provides the call logging middleware handler.
type CallLogger struct {
	logger zerolog.Logger
}

// CallLoggerParams defines the parameters for creating a call logging middleware.
type CallLoggerParams struct {
	fx.In

	Logger zerolog.Logger
}

// NewCallLogger creates a new call logging middleware instance.
func NewCallLogger(params CallLoggerParams) ICallLogger {
	return &CallLogger{
		logger: params.Logger,
	}
}

// Handler returns the call logging middleware function.
func (c *CallLogger) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.logger.Printf("path '%s' called by %s", r.URL.Path, r.RemoteAddr)
		h.ServeHTTP(w, r)
	})
}
