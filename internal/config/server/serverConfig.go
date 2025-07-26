package server

import (
	"time"

	config_util "github.com/greencoda/auth0-api-gateway/internal/util/config"
	"github.com/greencoda/confiq"
)

type Config struct {
	Address        string        `cfg:"address,default=:80"`
	ReadTimeout    time.Duration `cfg:"readTimeout,default=15s"`
	WriteTimeout   time.Duration `cfg:"writeTimeout,default=15s"`
	IdleTimeout    time.Duration `cfg:"idleTimeout,default=15s"`
	MaxHeaderBytes int           `cfg:"maxHeaderBytes,default=1048576"`
	ReleaseStage   string        `cfg:"releaseStage,default=local"`
	LogCalls       bool          `cfg:"logCalls,default=false"`
	LogLevel       string        `cfg:"logLevel,default=info"`
}

func NewConfig(configSet *confiq.ConfigSet) (*Config, error) {
	return config_util.LoadConfigFromSetWithPrefix[Config](configSet, "server")
}
