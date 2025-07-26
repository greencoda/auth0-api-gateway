package auth0

import (
	config_util "github.com/greencoda/auth0-api-gateway/internal/util/config"
	"github.com/greencoda/confiq"
)

type Config struct {
	Audience string `cfg:"audience,default=https://your-auth0-api.yourdomain.io"`
	Domain   string `cfg:"domain,default=your-auth0-tenant.eu.auth0.com"`
}

func NewConfig(configSet *confiq.ConfigSet) (*Config, error) {
	return config_util.LoadConfigFromSetWithPrefix[Config](configSet, "auth0")
}
