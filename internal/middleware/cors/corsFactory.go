package cors

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	subrouter_config "github.com/greencoda/auth0-api-gateway/internal/config/subrouter"
)

// Factory interface for creating CORS middleware for subrouters.
type ICORSFactory interface {
	NewCORS(config subrouter_config.CORSConfig) ICORS
}

// CORSFactory implements the ICORSFactory interface to create CORS middleware.
type CORSFactory struct{}

// NewCORS creates a new CORS middleware based on the provided configuration.
func NewCORSFactory() ICORSFactory {
	return &CORSFactory{}
}

// NewCORS creates a new CORS middleware based on the provided configuration.
func (c *CORSFactory) NewCORS(config subrouter_config.CORSConfig) ICORS {
	return &CORS{
		middlewareFunc: c.buildCORSMiddlewareFunc(config),
	}
}

func (c *CORSFactory) buildCORSMiddlewareFunc(config subrouter_config.CORSConfig) mux.MiddlewareFunc {
	corsOptions := []handlers.CORSOption{
		handlers.AllowedOrigins(config.AllowedOrigins),
		handlers.AllowedHeaders(config.AllowedHeaders),
		handlers.AllowedMethods(config.AllowedMethods),
		handlers.ExposedHeaders(config.ExposedHeaders),
		handlers.MaxAge(config.MaxAge),
	}

	if config.AllowCredentials {
		corsOptions = append(corsOptions, handlers.AllowCredentials())
	}

	return handlers.CORS(corsOptions...)
}
