package auth0_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	jwtvalidator "github.com/auth0/go-jwt-middleware/v2/validator"
	subrouter_config "github.com/greencoda/auth0-api-gateway/internal/config/subrouter"
	middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/auth0"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_Auth0ScopeValidator_Creation(t *testing.T) {
	Convey("When testing Auth0 scope validator creation", t, func() {
		auth0ValidatorFactory := middleware.NewAuth0ValidatorFactory()

		Convey("Should create validator with single valid scopes", func() {
			config := subrouter_config.AuthorizationConfig{
				RequiredScopes: []string{"read:all"},
			}

			validator := auth0ValidatorFactory.NewAuth0ScopeValidator(config)
			So(validator, ShouldNotBeNil)
			So(validator, ShouldImplement, (*middleware.IAuth0ScopeValidator)(nil))

			handler := validator.Handler()
			So(handler, ShouldNotBeNil)
		})

		Convey("Should create validator with multiple valid scopes", func() {
			config := subrouter_config.AuthorizationConfig{
				RequiredScopes: []string{"read:all", "write:users", "admin:system"},
			}

			validator := auth0ValidatorFactory.NewAuth0ScopeValidator(config)
			So(validator, ShouldNotBeNil)

			handler := validator.Handler()
			So(handler, ShouldNotBeNil)
		})

		Convey("Should create validator with no scopes", func() {
			config := subrouter_config.AuthorizationConfig{
				RequiredScopes: []string{},
			}

			validator := auth0ValidatorFactory.NewAuth0ScopeValidator(config)
			So(validator, ShouldNotBeNil)
		})

		Convey("Should create validator with nil scopes", func() {
			config := subrouter_config.AuthorizationConfig{
				RequiredScopes: nil,
			}

			validator := auth0ValidatorFactory.NewAuth0ScopeValidator(config)
			So(validator, ShouldNotBeNil)
		})
	})
}

func Test_Auth0ScopeValidator_Handler(t *testing.T) {
	Convey("When testing Auth0 scope validator handler functionality", t, func() {
		auth0ValidatorFactory := middleware.NewAuth0ValidatorFactory()

		Convey("Should allow access when token has all required scopes", func() {
			config := subrouter_config.AuthorizationConfig{
				RequiredScopes: []string{"read:all", "write:users"},
			}
			validator := auth0ValidatorFactory.NewAuth0ScopeValidator(config)
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			})

			middlewareHandler := validator.Handler()(testHandler)

			req := httptest.NewRequest("GET", "/test", nil)
			claims := &jwtvalidator.ValidatedClaims{
				CustomClaims: &middleware.CustomAuth0Claims{
					Scope: "read:all write:users admin:system",
				},
			}

			ctx := context.WithValue(req.Context(), jwtmiddleware.ContextKey{}, claims)
			req = req.WithContext(ctx)

			recorder := httptest.NewRecorder()

			middlewareHandler.ServeHTTP(recorder, req)

			So(recorder.Code, ShouldEqual, http.StatusOK)
			So(recorder.Body.String(), ShouldEqual, "success")
		})

		Convey("Should deny access when token missing required scopes", func() {
			config := subrouter_config.AuthorizationConfig{
				RequiredScopes: []string{"read:all", "admin:system"},
			}

			validator := auth0ValidatorFactory.NewAuth0ScopeValidator(config)

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			})

			middlewareHandler := validator.Handler()(testHandler)

			req := httptest.NewRequest("GET", "/test", nil)

			claims := &jwtvalidator.ValidatedClaims{
				CustomClaims: &middleware.CustomAuth0Claims{
					Scope: "read:all write:users", // missing admin:system
				},
			}

			ctx := context.WithValue(req.Context(), jwtmiddleware.ContextKey{}, claims)
			req = req.WithContext(ctx)

			recorder := httptest.NewRecorder()

			middlewareHandler.ServeHTTP(recorder, req)

			So(recorder.Code, ShouldEqual, http.StatusForbidden)
			So(recorder.Header().Get("Content-Type"), ShouldEqual, "application/json")
			So(recorder.Body.String(), ShouldEqual, `{"message":"Insufficient access privileges."}`)
		})

		Convey("Should deny access when no token in context", func() {
			config := subrouter_config.AuthorizationConfig{
				RequiredScopes: []string{"read:all"},
			}

			validator := auth0ValidatorFactory.NewAuth0ScopeValidator(config)

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			})

			middlewareHandler := validator.Handler()(testHandler)

			req := httptest.NewRequest("GET", "/test", nil)

			recorder := httptest.NewRecorder()

			middlewareHandler.ServeHTTP(recorder, req)

			So(recorder.Code, ShouldEqual, http.StatusBadRequest)
			So(recorder.Header().Get("Content-Type"), ShouldEqual, "application/json")
			So(recorder.Body.String(), ShouldEqual, `{"message":"Cannot access token."}`)
		})

		Convey("Should deny access when token has invalid custom claims", func() {
			config := subrouter_config.AuthorizationConfig{
				RequiredScopes: []string{"read:all"},
			}

			validator := auth0ValidatorFactory.NewAuth0ScopeValidator(config)

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			})

			middlewareHandler := validator.Handler()(testHandler)

			req := httptest.NewRequest("GET", "/test", nil)

			claims := &jwtvalidator.ValidatedClaims{
				CustomClaims: nil,
			}

			ctx := context.WithValue(req.Context(), jwtmiddleware.ContextKey{}, claims)
			req = req.WithContext(ctx)

			recorder := httptest.NewRecorder()

			middlewareHandler.ServeHTTP(recorder, req)

			So(recorder.Code, ShouldEqual, http.StatusBadRequest)
			So(recorder.Header().Get("Content-Type"), ShouldEqual, "application/json")
			So(recorder.Body.String(), ShouldEqual, `{"message":"Invalid claims in token."}`)
		})

		Convey("Should allow access when no scopes required", func() {
			config := subrouter_config.AuthorizationConfig{
				RequiredScopes: []string{},
			}

			validator := auth0ValidatorFactory.NewAuth0ScopeValidator(config)

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			})

			middlewareHandler := validator.Handler()(testHandler)

			req := httptest.NewRequest("GET", "/test", nil)

			claims := &jwtvalidator.ValidatedClaims{
				CustomClaims: &middleware.CustomAuth0Claims{
					Scope: "any:scope",
				},
			}

			ctx := context.WithValue(req.Context(), jwtmiddleware.ContextKey{}, claims)
			req = req.WithContext(ctx)

			recorder := httptest.NewRecorder()

			middlewareHandler.ServeHTTP(recorder, req)

			So(recorder.Code, ShouldEqual, http.StatusOK)
			So(recorder.Body.String(), ShouldEqual, "success")
		})

		Convey("Should allow access when nil scopes required", func() {
			config := subrouter_config.AuthorizationConfig{
				RequiredScopes: nil,
			}

			validator := auth0ValidatorFactory.NewAuth0ScopeValidator(config)

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			})

			middlewareHandler := validator.Handler()(testHandler)

			req := httptest.NewRequest("GET", "/test", nil)

			claims := &jwtvalidator.ValidatedClaims{
				CustomClaims: &middleware.CustomAuth0Claims{
					Scope: "any:scope",
				},
			}

			ctx := context.WithValue(req.Context(), jwtmiddleware.ContextKey{}, claims)
			req = req.WithContext(ctx)

			recorder := httptest.NewRecorder()

			middlewareHandler.ServeHTTP(recorder, req)

			So(recorder.Code, ShouldEqual, http.StatusOK)
			So(recorder.Body.String(), ShouldEqual, "success")
		})
	})
}
