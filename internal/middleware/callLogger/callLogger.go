package callLogger

import (
	"net/http"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

type ICallLogger interface {
	Handler(h http.Handler) http.Handler
}

type CallLogger struct {
	logger zerolog.Logger
}

type CallLoggerParams struct {
	fx.In

	Logger zerolog.Logger
}

func NewCallLogger(params CallLoggerParams) ICallLogger {
	return &CallLogger{
		logger: params.Logger,
	}
}

// Handler returns the call logging middleware function.
func (c *CallLogger) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
		c.logger.Printf("path '%s' called by %s", req.URL.Path, req.RemoteAddr)
		h.ServeHTTP(responseWriter, req)
	})
}
