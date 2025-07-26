package subrouter

import (
	"time"

	config_util "github.com/greencoda/auth0-api-gateway/internal/util/config"
	"github.com/greencoda/confiq"
)

type AuthorizationConfig struct {
	RequiredScopes []string `cfg:"requiredScopes"`
}

type RateLimitConfig struct {
	Limit              int64         `cfg:"maxRequests"`
	Period             time.Duration `cfg:"expiration"`
	TrustForwardHeader bool          `cfg:"trustForwardHeader,default=false"`
}

type CORSConfig struct {
	AllowedOrigins     []string `cfg:"allowedOrigins"`
	AllowedMethods     []string `cfg:"allowedMethods"`
	AllowedHeaders     []string `cfg:"allowedHeaders"`
	ExposedHeaders     []string `cfg:"exposedHeaders"`
	AllowCredentials   bool     `cfg:"allowCredentials"`
	MaxAge             int      `cfg:"maxAge"`
	OptionsPassthrough bool     `cfg:"optionsPassthrough"`
	Debug              bool     `cfg:"debug"`
}

type SubrouterConfig struct {
	Name                string               `cfg:"name"`
	TargetURL           string               `cfg:"targetURL"`
	Prefix              string               `cfg:"prefix"`
	StripPrefix         bool                 `cfg:"stripPrefix,default=false"`
	AuthorizationConfig *AuthorizationConfig `cfg:"authorizationConfig"`
	RateLimitConfig     *RateLimitConfig     `cfg:"rateLimit"`
	GZip                bool                 `cfg:"gzip,default=false"`
	CORSConfig          *CORSConfig          `cfg:"corsConfig"`
}

type Config []SubrouterConfig

func NewConfig(configSet *confiq.ConfigSet) (*Config, error) {
	return config_util.LoadConfigFromSetWithPrefix[Config](configSet, "subrouters")
}
