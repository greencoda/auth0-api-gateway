package server_test

import (
	"testing"

	mock_auth0_middleware "github.com/greencoda/auth0-api-gateway/internal/mocks/middleware/auth0"
	mock_callLogger_middleware "github.com/greencoda/auth0-api-gateway/internal/mocks/middleware/callLogger"
	mock_cors_middleware "github.com/greencoda/auth0-api-gateway/internal/mocks/middleware/cors"
	mock_rateLimit_middleware "github.com/greencoda/auth0-api-gateway/internal/mocks/middleware/rateLimit"
	"github.com/greencoda/auth0-api-gateway/internal/server"
	"github.com/rs/zerolog"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewServer(t *testing.T) {
	Convey("When creating a new server", t, func() {
		var (
			testLogger = zerolog.New(zerolog.NewConsoleWriter())

			// Mocks
			mockAuth0ValidatorFactory mock_auth0_middleware.IAuth0ValidatorFactory
			mockAuth0TokenValidator   mock_auth0_middleware.IAuth0TokenValidator
			mockAuth0ScopeValidator   mock_auth0_middleware.IAuth0ScopeValidator
			mockCORSFactory           mock_cors_middleware.ICORSFactory
			mockICORS                 mock_cors_middleware.ICORS
			mockRateLimitFactory      mock_rateLimit_middleware.IRateLimitFactory
			mockRateLimit             mock_rateLimit_middleware.IRateLimit
			mockCallLogger            mock_callLogger_middleware.ICallLogger
		)

		Convey("With fully valid config", func() {
			Convey("When Auth0 Token and Scope validator can be set up", func() {
				mockAuth0ValidatorFactory.On("NewAuth0TokenValidator", validAuth0Config).Return(&mockAuth0TokenValidator, nil)
				mockAuth0ValidatorFactory.On("NewAuth0ScopeValidator", *validSubrouterConfigs[0].AuthorizationConfig).Return(&mockAuth0ScopeValidator)
				mockRateLimitFactory.On("NewRateLimit", *(validSubrouterConfigs[0].RateLimitConfig)).Return(&mockRateLimit, nil)
				mockCORSFactory.On("NewCORS", *validSubrouterConfigs[0].CORSConfig).Return(&mockICORS, nil)

				mockAuth0TokenValidator.On("Handler").Return(noopMiddlewareFunc)
				mockAuth0ScopeValidator.On("Handler").Return(noopMiddlewareFunc)
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
						CallLogMiddleware:          &mockCallLogger,
						Logger:                     testLogger,
					},
				)
				So(reverseProxyHandler, ShouldNotBeNil)
				So(err, ShouldBeNil)

				testServer, err := server.NewServer(
					server.ServerParams{
						ServerConfig:        &validServerConfig,
						ReverseProxyHandler: reverseProxyHandler,
						Logger:              testLogger,
					},
				)
				So(testServer, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})
		})
	})
}
