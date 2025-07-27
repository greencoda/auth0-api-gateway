package auth0

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gorilla/mux"
	subrouter_config "github.com/greencoda/auth0-api-gateway/internal/config/subrouter"
)

type IAuth0ScopeValidator interface {
	Handler() mux.MiddlewareFunc
}

type Auth0ScopeValidator struct {
	middlewareFunc mux.MiddlewareFunc
}

func (a *Auth0ScopeValidator) Handler() mux.MiddlewareFunc {
	return a.middlewareFunc
}

func buildAuth0ScopeMiddlewareFunc(config subrouter_config.AuthorizationConfig) mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
			token, isValidatedClaim := req.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
			if !isValidatedClaim {
				handleScopeError(responseWriter, http.StatusBadRequest, "Cannot access token.")

				return
			}

			customAuth0Claims, isValidatedClaim := token.CustomClaims.(*CustomAuth0Claims)
			if !isValidatedClaim {
				handleScopeError(responseWriter, http.StatusBadRequest, "Invalid claims in token.")

				return
			}

			if !customAuth0Claims.HasAllScopes(config.RequiredScopes) {
				handleScopeError(responseWriter, http.StatusForbidden, "Insufficient access privileges.")

				return
			}

			handler.ServeHTTP(responseWriter, req)
		})
	}
}

func handleScopeError(responseWriter http.ResponseWriter, httpStatusCode int, message string) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(httpStatusCode)
	_, _ = responseWriter.Write([]byte(`{"message":"` + message + `"}`))
}
