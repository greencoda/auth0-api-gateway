package server_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	auth0_config "github.com/greencoda/auth0-api-gateway/internal/config/auth0"
	server_config "github.com/greencoda/auth0-api-gateway/internal/config/server"
	subrouter_config "github.com/greencoda/auth0-api-gateway/internal/config/subrouter"
	mock_auth0_middleware "github.com/greencoda/auth0-api-gateway/internal/mocks/middleware/auth0"
	mock_cors_middleware "github.com/greencoda/auth0-api-gateway/internal/mocks/middleware/cors"
	mock_rateLimit_middleware "github.com/greencoda/auth0-api-gateway/internal/mocks/middleware/rateLimit"
	mock_requestLogger_middleware "github.com/greencoda/auth0-api-gateway/internal/mocks/middleware/requestLogger"
	"github.com/greencoda/auth0-api-gateway/internal/server"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	validServerConfig = server_config.Config{
		Address:        ":88",
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1048576,
		ReleaseStage:   "local",
		LogRequests:    true,
		LogLevel:       "info",
	}
	validAuth0Config = auth0_config.Config{
		Audience: "",
		Domain:   "",
	}
	validSubrouterConfigs = subrouter_config.Config{
		{
			Name:        "Test API",
			TargetURL:   "http://localhost:8088",
			Prefix:      "/protected",
			StripPrefix: true,
			AuthorizationConfig: &subrouter_config.AuthorizationConfig{
				RequiredScopes: []string{"read:all"},
			},
			RateLimitConfig: &subrouter_config.RateLimitConfig{
				Limit:              5,
				Period:             time.Second,
				TrustForwardHeader: true,
			},
			GZip: true,
			CORSConfig: &subrouter_config.CORSConfig{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET"},
				AllowedHeaders:   []string{"Authorization"},
				AllowCredentials: true,
				MaxAge:           86400,
			},
		},
	}

	invalidSubrouterConfigs = subrouter_config.Config{
		{
			Name:                "Test API",
			TargetURL:           " $$$$:// ",
			Prefix:              "/protected",
			StripPrefix:         true,
			AuthorizationConfig: nil,
			RateLimitConfig:     nil,
			GZip:                true,
			CORSConfig:          nil,
		},
	}
	noopMiddlewareFunc mux.MiddlewareFunc = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
		})
	}

	testError = errors.New("test error")
)

func Test_NewReverseProxyHandler(t *testing.T) {
	Convey("When creating a new server", t, func() {
		var (
			testLogger = zerolog.New(zerolog.NewConsoleWriter())

			mockAuth0ValidatorFactory mock_auth0_middleware.IAuth0ValidatorFactory
			mockAuth0TokenValidator   mock_auth0_middleware.IAuth0TokenValidator
			mockCORSFactory           mock_cors_middleware.ICORSFactory
			mockICORS                 mock_cors_middleware.ICORS
			mockRateLimitFactory      mock_rateLimit_middleware.IRateLimitFactory
			mockRateLimit             mock_rateLimit_middleware.IRateLimit
			mockRequestLogger         mock_requestLogger_middleware.IRequestLogger
		)

		Convey("With fully valid config", func() {
			Convey("When Auth0 Token validator cannot be set up", func() {
				mockAuth0ValidatorFactory.On("NewAuth0TokenValidator", validAuth0Config).Return(nil, testError)

				reverseProxyHandler, err := server.NewReverseProxyHandler(
					server.ReverseProxyHandlerParams{
						Auth0Config:                &validAuth0Config,
						ServerConfig:               &validServerConfig,
						SubrouterConfigs:           &validSubrouterConfigs,
						Auth0MiddlewareFactory:     &mockAuth0ValidatorFactory,
						CORSMiddlewareFactory:      &mockCORSFactory,
						RateLimitMiddlewareFactory: &mockRateLimitFactory,
						RequestLoggerMiddleware:    &mockRequestLogger,
						Logger:                     testLogger,
					},
				)
				So(reverseProxyHandler, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})

			Convey("When Auth0 Scope validator cannot be set up", func() {
				mockAuth0ValidatorFactory.On("NewAuth0TokenValidator", validAuth0Config).Return(&mockAuth0TokenValidator, nil)
				mockAuth0ValidatorFactory.On("NewAuth0ScopeValidator", *validSubrouterConfigs[0].AuthorizationConfig).Return(nil)
				mockRateLimitFactory.On("NewRateLimit", *(validSubrouterConfigs[0].RateLimitConfig)).Return(&mockRateLimit, nil)
				mockCORSFactory.On("NewCORS", *validSubrouterConfigs[0].CORSConfig).Return(&mockICORS, nil)

				mockAuth0TokenValidator.On("Handler").Return(noopMiddlewareFunc)
				mockRateLimit.On("Handler").Return(noopMiddlewareFunc)
				mockICORS.On("Handler").Return(noopMiddlewareFunc)

				reverseProxyHandler, err := server.NewReverseProxyHandler(
					server.ReverseProxyHandlerParams{
						Auth0Config:                &validAuth0Config,
						ServerConfig:               &validServerConfig,
						SubrouterConfigs:           &validSubrouterConfigs,
						Auth0MiddlewareFactory:     &mockAuth0ValidatorFactory,
						CORSMiddlewareFactory:      &mockCORSFactory,
						RateLimitMiddlewareFactory: &mockRateLimitFactory,
						RequestLoggerMiddleware:    &mockRequestLogger,
						Logger:                     testLogger,
					},
				)
				So(reverseProxyHandler, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("With invalid target URL in config", func() {
			mockAuth0ValidatorFactory.On("NewAuth0TokenValidator", validAuth0Config).Return(&mockAuth0TokenValidator, nil)

			reverseProxyHandler, err := server.NewReverseProxyHandler(
				server.ReverseProxyHandlerParams{
					Auth0Config:                &validAuth0Config,
					ServerConfig:               &validServerConfig,
					SubrouterConfigs:           &invalidSubrouterConfigs,
					Auth0MiddlewareFactory:     &mockAuth0ValidatorFactory,
					CORSMiddlewareFactory:      &mockCORSFactory,
					RateLimitMiddlewareFactory: &mockRateLimitFactory,
					RequestLoggerMiddleware:    &mockRequestLogger,
					Logger:                     testLogger,
				},
			)
			So(reverseProxyHandler, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}
