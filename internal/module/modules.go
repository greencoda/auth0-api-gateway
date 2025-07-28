package module

import (
	auth0_config "github.com/greencoda/auth0-api-gateway/internal/config/auth0"
	server_config "github.com/greencoda/auth0-api-gateway/internal/config/server"
	subrouter_config "github.com/greencoda/auth0-api-gateway/internal/config/subrouter"
	auth0_middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/auth0"
	cors_middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/cors"
	rateLimit_middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/rateLimit"
	requestLogger_middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/requestLogger"
	"github.com/greencoda/auth0-api-gateway/internal/server"
	config_util "github.com/greencoda/auth0-api-gateway/internal/util/config"
	logging_util "github.com/greencoda/auth0-api-gateway/internal/util/logging"
	"go.uber.org/fx"
)

var apiConfigModule = fx.Module(
	"config",
	fx.Provide(
		config_util.LoadConfigYAML,
		auth0_config.NewConfig,
		server_config.NewConfig,
		subrouter_config.NewConfig,
	),
)

var logicModule = fx.Module(
	"logic",
	fx.Provide(
		requestLogger_middleware.NewMiddleware,
		auth0_middleware.NewAuth0ValidatorFactory,
		cors_middleware.NewCORSFactory,
		rateLimit_middleware.NewRateLimitFactory,
		server.NewReverseProxyHandler,
		server.NewServer,
	),
)

var observabilityModule = fx.Module(
	"observability",
	fx.Provide(
		logging_util.NewZeroLogger,
	),
)

// NewServerModule creates a new server module with the specified config filename
func NewServerModule(configFilename string) *fx.App {
	return fx.New(
		fx.Supply(
			config_util.ConfigFilename(configFilename),
		),
		apiConfigModule,
		fx.WithLogger(
			logging_util.NewFXLogger,
		),
		logicModule,
		observabilityModule,
		fx.Invoke(Launcher),
	)
}
