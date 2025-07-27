package auth0

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gorilla/mux"
	auth0_config "github.com/greencoda/auth0-api-gateway/internal/config/auth0"
)

const jwtCacheTTL = time.Duration(5 * time.Minute)

type IAuth0TokenValidator interface {
	Handler() mux.MiddlewareFunc
}

type Auth0TokenValidator struct {
	middlewareFunc mux.MiddlewareFunc
}

func (a *Auth0TokenValidator) Handler() mux.MiddlewareFunc {
	return a.middlewareFunc
}

func buildJWTMiddlewareFunc(config auth0_config.Config) (mux.MiddlewareFunc, error) {
	issuerURL, err := url.Parse("https://" + config.Domain + "/")
	if err != nil {
		return nil, fmt.Errorf("failed to parse the issuer url: %w", err)
	}

	var (
		audiences       []string
		cachingProvider = jwks.NewCachingProvider(issuerURL, jwtCacheTTL)
	)

	if config.Audience != "" {
		audiences = []string{config.Audience}
	}

	jwtValidator, err := validator.New(
		cachingProvider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		audiences,
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomAuth0Claims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to set up the jwt validator: %w", err)
	}

	errorHandler := func(responseWriter http.ResponseWriter, req *http.Request, err error) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusUnauthorized)
		_, _ = responseWriter.Write([]byte(`{"message":"Failed to validate JWT: ` + err.Error() + `"}`))
	}

	return jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	).CheckJWT, nil
}
