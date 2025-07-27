package module

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

type LauncherParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Server    *http.Server
	Logger    zerolog.Logger
}

func Launcher(params LauncherParams) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			params.Logger.Print("Starting Auth0 API Gateway")
			go func() {
				err := params.Server.ListenAndServe()
				if err != nil {
					params.Logger.Fatal().Err(err).Msg("Failed to start server")
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			params.Logger.Print("Shutting down Auth0 API Gateway")

			return nil
		},
	})
}
