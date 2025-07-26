package auth0

import (
	"context"
	"strings"
)

type ICustomAuth0Claims interface {
	Validate(context.Context) error
	HasAllScopes([]string) bool
}

type CustomAuth0Claims struct {
	Scope string `json:"scope"`
}

func (c CustomAuth0Claims) Validate(ctx context.Context) error {
	return nil
}

func (c CustomAuth0Claims) HasAllScopes(expectedScopes []string) bool {
	scopes := strings.Split(c.Scope, " ")
	for _, expectedScope := range expectedScopes {
		if !scopeInSlice(expectedScope, scopes) {
			return false
		}
	}

	return true
}

func scopeInSlice(expectedScope string, scopes []string) bool {
	for _, scope := range scopes {
		if scope == expectedScope {
			return true
		}
	}

	return false
}
