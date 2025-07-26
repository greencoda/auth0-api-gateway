package server

import (
	"log"
	"net/http"

	server_config "github.com/greencoda/auth0-api-gateway/internal/config/server"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

type ServerParams struct {
	fx.In

	ServerConfig        *server_config.Config
	ReverseProxyHandler IReverseProxyHandler
	Logger              zerolog.Logger
}

func NewServer(params ServerParams) (*http.Server, error) {
	stdLogger := log.New(params.Logger, "", 0)

	return &http.Server{
		Handler:        params.ReverseProxyHandler,
		Addr:           params.ServerConfig.Address,
		ReadTimeout:    params.ServerConfig.ReadTimeout,
		WriteTimeout:   params.ServerConfig.WriteTimeout,
		IdleTimeout:    params.ServerConfig.IdleTimeout,
		MaxHeaderBytes: params.ServerConfig.MaxHeaderBytes,
		ErrorLog:       stdLogger,
	}, nil
}
