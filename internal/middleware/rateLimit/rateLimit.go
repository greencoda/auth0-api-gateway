package rateLimit

import (
	"github.com/gorilla/mux"
	subrouter_config "github.com/greencoda/auth0-api-gateway/internal/config/subrouter"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

type IRateLimitFactory interface {
	NewRateLimit(config subrouter_config.RateLimitConfig) IRateLimit
}

type RateLimitFactory struct{}

func (r *RateLimitFactory) NewRateLimit(config subrouter_config.RateLimitConfig) IRateLimit {
	return &RateLimit{
		middlewareFunc: buildRateLimiterFunc(config),
	}
}

func NewRateLimitFactory() IRateLimitFactory {
	return &RateLimitFactory{}
}

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
	rate := limiter.Rate{
		Period: config.Period,
		Limit:  config.Limit,
	}

	rateLimiter := limiter.New(
		memory.NewStore(),
		rate,
		limiter.WithTrustForwardHeader(config.TrustForwardHeader),
	)

	middleware := stdlib.NewMiddleware(
		rateLimiter,
	)

	return middleware.Handler
}
