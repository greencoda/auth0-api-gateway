package server

import (
	"net/http"
	"net/url"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	auth0_config "github.com/greencoda/auth0-api-gateway/internal/config/auth0"
	server_config "github.com/greencoda/auth0-api-gateway/internal/config/server"
	subrouter_config "github.com/greencoda/auth0-api-gateway/internal/config/subrouter"
	auth0_middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/auth0"
	callLogger_middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/callLogger"
	cors_middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/cors"
	rateLimit_middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/rateLimit"
	reverseProxy_util "github.com/greencoda/auth0-api-gateway/internal/util/reverseProxy"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

type IReverseProxyHandler http.Handler

type ReverseProxyHandlerParams struct {
	fx.In

	Auth0Config      *auth0_config.Config
	ServerConfig     *server_config.Config
	SubrouterConfigs *subrouter_config.Config

	Auth0MiddlewareFactory     auth0_middleware.IAuth0ValidatorFactory
	CallLogMiddleware          callLogger_middleware.ICallLogger
	CORSMiddlewareFactory      cors_middleware.ICORSFactory
	RateLimitMiddlewareFactory rateLimit_middleware.IRateLimitFactory

	Logger zerolog.Logger
}

func NewReverseProxyHandler(params ReverseProxyHandlerParams) (IReverseProxyHandler, error) {
	router := mux.NewRouter()

	if params.ServerConfig.LogCalls {
		router.Use(params.CallLogMiddleware.Handler)
		params.Logger.Print("Logging API calls is enabled")
	}

	auth0TokenValidatorMiddleware, err := params.Auth0MiddlewareFactory.NewAuth0TokenValidator(*params.Auth0Config)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to set up Auth0 token validator middleware")
	}

	for _, subrouterConfig := range *params.SubrouterConfigs {
		subRouter := router.PathPrefix(subrouterConfig.Prefix).Subrouter()

		targetURL, err := url.Parse(subrouterConfig.TargetURL)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse target API URL '%s' of subrouter '%s'", subrouterConfig.TargetURL, subrouterConfig.Name)
		}

		if subrouterConfig.RateLimitConfig != nil {
			rateLimiterMiddleware := params.RateLimitMiddlewareFactory.NewRateLimit(*subrouterConfig.RateLimitConfig)
			subRouter.Use(rateLimiterMiddleware.Handler())
		}

		if subrouterConfig.CORSConfig != nil {
			corsMiddleware := params.CORSMiddlewareFactory.NewCORS(*subrouterConfig.CORSConfig)
			subRouter.Use(corsMiddleware.Handler())
		}

		if subrouterConfig.AuthorizationConfig != nil {
			subRouter.Use(auth0TokenValidatorMiddleware.Handler())

			if len(subrouterConfig.AuthorizationConfig.RequiredScopes) > 0 {
				auth0ScopeValidatorMiddleware := params.Auth0MiddlewareFactory.NewAuth0ScopeValidator(*subrouterConfig.AuthorizationConfig)
				if auth0ScopeValidatorMiddleware == nil {
					return nil, errors.New("failed to create Auth0 scope validator middleware")
				}

				subRouter.Use(auth0ScopeValidatorMiddleware.Handler())
			}
		}

		if subrouterConfig.GZip {
			subRouter.Use(handlers.CompressHandler)
		}

		reverseProxy := reverseProxy_util.NewReverseProxy(targetURL)

		var subRouterHandler http.Handler = reverseProxy
		if subrouterConfig.StripPrefix {
			subRouterHandler = http.StripPrefix(subrouterConfig.Prefix, subRouterHandler)
		}

		subRouter.NewRoute().Handler(
			subRouterHandler,
		)

		params.Logger.Printf("API router for '%s' setup complete on path prefix: %s", subrouterConfig.Name, subrouterConfig.Prefix)
	}

	return router, nil
}
