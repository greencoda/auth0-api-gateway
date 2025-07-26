package auth0

import (
	auth0_config "github.com/greencoda/auth0-api-gateway/internal/config/auth0"
	subrouter_config "github.com/greencoda/auth0-api-gateway/internal/config/subrouter"
)

type IAuth0ValidatorFactory interface {
	NewAuth0ScopeValidator(config subrouter_config.AuthorizationConfig) IAuth0ScopeValidator
	NewAuth0TokenValidator(config auth0_config.Config) (IAuth0TokenValidator, error)
}

func NewAuth0ValidatorFactory() IAuth0ValidatorFactory {
	return &Auth0ValidatorFactory{}
}

type Auth0ValidatorFactory struct{}

func (a *Auth0ValidatorFactory) NewAuth0ScopeValidator(config subrouter_config.AuthorizationConfig) IAuth0ScopeValidator {
	return &Auth0ScopeValidator{
		middlewareFunc: buildAuth0ScopeMiddlewareFunc(config),
	}
}

func (a *Auth0ValidatorFactory) NewAuth0TokenValidator(config auth0_config.Config) (IAuth0TokenValidator, error) {
	jwtMiddlewareFunc, err := buildJWTMiddlewareFunc(config)
	if err != nil {
		return nil, err
	}

	return &Auth0TokenValidator{
		middlewareFunc: jwtMiddlewareFunc,
	}, nil
}
