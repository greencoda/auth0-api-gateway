package rateLimit

import (
	subrouter_config "github.com/greencoda/auth0-api-gateway/internal/config/subrouter"
)

type IRateLimitFactory interface {
	NewRateLimit(config subrouter_config.RateLimitConfig) IRateLimit
}

type RateLimitFactory struct{}

func NewRateLimitFactory() IRateLimitFactory {
	return &RateLimitFactory{}
}

func (r *RateLimitFactory) NewRateLimit(config subrouter_config.RateLimitConfig) IRateLimit {
	return &RateLimit{
		middlewareFunc: buildRateLimiterFunc(config),
	}
}
