package module

import (
	auth0_config "github.com/greencoda/auth0-api-gateway/internal/config/auth0"
	server_config "github.com/greencoda/auth0-api-gateway/internal/config/server"
	subrouter_config "github.com/greencoda/auth0-api-gateway/internal/config/subrouter"
	auth0_middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/auth0"
	callLogger_middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/callLogger"
	cors_middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/cors"
	rateLimit_middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/rateLimit"
	"github.com/greencoda/auth0-api-gateway/internal/server"
	config_util "github.com/greencoda/auth0-api-gateway/internal/util/config"
	logging_util "github.com/greencoda/auth0-api-gateway/internal/util/logging"
	"go.uber.org/fx"
)

const apiConfigFilename = "config.yaml"

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
	"business",
	fx.Provide(
		auth0_middleware.NewAuth0ValidatorFactory,
		callLogger_middleware.NewCallLogger,
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

var ServerModule = fx.New(
	fx.Supply(
		config_util.ConfigFilename(apiConfigFilename),
	),
	apiConfigModule,
	fx.WithLogger(
		logging_util.NewFXLogger,
	),
	logicModule,
	observabilityModule,
	fx.Invoke(Launcher),
)
