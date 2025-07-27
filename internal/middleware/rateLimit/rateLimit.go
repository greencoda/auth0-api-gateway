package rateLimit

import (
	"github.com/gorilla/mux"
	subrouter_config "github.com/greencoda/auth0-api-gateway/internal/config/subrouter"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

type IRateLimit interface {
	Handler() mux.MiddlewareFunc
}

type RateLimit struct {
	middlewareFunc mux.MiddlewareFunc
}

func (c *RateLimit) Handler() mux.MiddlewareFunc {
	return c.middlewareFunc
}

func buildRateLimiterFunc(config subrouter_config.RateLimitConfig) mux.MiddlewareFunc {
	var (
		limiterStore = memory.NewStore()
		limiterRate  = limiter.Rate{
			Period: config.Period,
			Limit:  config.Limit,
		}
	)

	middleware := stdlib.NewMiddleware(
		limiter.New(
			limiterStore,
			limiterRate,
			limiter.WithTrustForwardHeader(config.TrustForwardHeader),
		),
	)

	return middleware.Handler
}
