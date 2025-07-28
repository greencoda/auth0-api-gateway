package requestLogger

import (
	"net/http"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

type IRequestLogger interface {
	Handler(h http.Handler) http.Handler
}

type RequestLogger struct {
	logger zerolog.Logger
}

type RequestLoggerParams struct {
	fx.In

	Logger zerolog.Logger
}

func NewMiddleware(params RequestLoggerParams) IRequestLogger {
	return &RequestLogger{
		logger: params.Logger,
	}
}

func (c *RequestLogger) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
		c.logger.Info().
			Str("method", req.Method).
			Str("path", req.URL.Path).
			Str("remote_addr", req.RemoteAddr).
			Str("user_agent", req.UserAgent()).
			Msg("Request received")

		h.ServeHTTP(responseWriter, req)
	})
}
