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
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
			if !ok {
				handleScopeError(w, http.StatusBadRequest, "Cannot access token.")
				return
			}

			customAuth0Claims, ok := token.CustomClaims.(*CustomAuth0Claims)
			if !ok {
				handleScopeError(w, http.StatusBadRequest, "Invalid claims in token.")
				return
			}

			if !customAuth0Claims.HasAllScopes(config.RequiredScopes) {
				handleScopeError(w, http.StatusForbidden, "Insufficient access privileges.")
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

func handleScopeError(w http.ResponseWriter, httpStatusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	_, _ = w.Write([]byte(`{"message":"` + message + `"}`))
}
