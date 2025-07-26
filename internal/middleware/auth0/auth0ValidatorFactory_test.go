package auth0_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	jwtvalidator "github.com/auth0/go-jwt-middleware/v2/validator"
	auth0_config "github.com/greencoda/auth0-api-gateway/internal/config/auth0"
	subrouter_config "github.com/greencoda/auth0-api-gateway/internal/config/subrouter"
	middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/auth0"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewAuth0ValidatorFactory(t *testing.T) {
	Convey("When creating a new Auth0 validator factory", t, func() {
		factory := middleware.NewAuth0ValidatorFactory()
		So(factory, ShouldNotBeNil)
		So(factory, ShouldImplement, (*middleware.IAuth0ValidatorFactory)(nil))
	})
}

func Test_Auth0ValidatorFactory_NewAuth0ScopeValidator(t *testing.T) {
	Convey("When creating Auth0 scope validator", t, func() {
		factory := middleware.NewAuth0ValidatorFactory()

		Convey("With valid authorization config", func() {
			config := subrouter_config.AuthorizationConfig{
				RequiredScopes: []string{"read:all", "write:users"},
			}

			validator := factory.NewAuth0ScopeValidator(config)
			So(validator, ShouldNotBeNil)
			So(validator, ShouldImplement, (*middleware.IAuth0ScopeValidator)(nil))

			handler := validator.Handler()
			So(handler, ShouldNotBeNil)
		})

		Convey("With empty scopes", func() {
			config := subrouter_config.AuthorizationConfig{
				RequiredScopes: []string{},
			}

			validator := factory.NewAuth0ScopeValidator(config)
			So(validator, ShouldNotBeNil)

			handler := validator.Handler()
			So(handler, ShouldNotBeNil)
		})

		Convey("With nil scopes", func() {
			config := subrouter_config.AuthorizationConfig{
				RequiredScopes: nil,
			}

			validator := factory.NewAuth0ScopeValidator(config)
			So(validator, ShouldNotBeNil)

			handler := validator.Handler()
			So(handler, ShouldNotBeNil)
		})

		Convey("With single scope", func() {
			config := subrouter_config.AuthorizationConfig{
				RequiredScopes: []string{"admin:access"},
			}

			validator := factory.NewAuth0ScopeValidator(config)
			So(validator, ShouldNotBeNil)

			handler := validator.Handler()
			So(handler, ShouldNotBeNil)
		})

		Convey("With multiple scopes", func() {
			config := subrouter_config.AuthorizationConfig{
				RequiredScopes: []string{"read:users", "write:users", "delete:users", "admin:system"},
			}

			validator := factory.NewAuth0ScopeValidator(config)
			So(validator, ShouldNotBeNil)

			handler := validator.Handler()
			So(handler, ShouldNotBeNil)
		})
	})
}

func Test_Auth0ValidatorFactory_NewAuth0TokenValidator(t *testing.T) {
	Convey("When creating Auth0 token validator", t, func() {
		factory := middleware.NewAuth0ValidatorFactory()

		Convey("With valid Auth0 config", func() {
			config := auth0_config.Config{
				Audience: "https://api.example.com",
				Domain:   "example.auth0.com",
			}

			validator, err := factory.NewAuth0TokenValidator(config)
			So(err, ShouldBeNil)
			So(validator, ShouldNotBeNil)
			So(validator, ShouldImplement, (*middleware.IAuth0TokenValidator)(nil))
		})

		Convey("With invalid domain", func() {
			config := auth0_config.Config{
				Audience: "https://api.example.com",
				Domain:   "invalid domain with spaces",
			}

			validator, err := factory.NewAuth0TokenValidator(config)
			So(err, ShouldNotBeNil)
			So(validator, ShouldBeNil)
		})

		Convey("With empty config", func() {
			config := auth0_config.Config{}

			validator, err := factory.NewAuth0TokenValidator(config)
			So(err, ShouldNotBeNil)
			So(validator, ShouldBeNil)
		})

		Convey("With minimal valid config", func() {
			config := auth0_config.Config{
				Audience: "test-audience",
				Domain:   "test.auth0.com",
			}

			validator, err := factory.NewAuth0TokenValidator(config)
			So(err, ShouldBeNil)
			So(validator, ShouldNotBeNil)
		})
	})
}

func Test_Auth0Middleware_ComprehensiveCoverage(t *testing.T) {
	Convey("When testing comprehensive Auth0 middleware coverage", t, func() {
		factory := middleware.NewAuth0ValidatorFactory()

		Convey("Should handle all middleware code paths", func() {
			config := subrouter_config.AuthorizationConfig{
				RequiredScopes: []string{"admin:read"},
			}

			validator := factory.NewAuth0ScopeValidator(config)
			So(validator, ShouldNotBeNil)

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("success"))
			})

			middlewareHandler := validator.Handler()(testHandler)

			Convey("Test with completely wrong context type", func() {
				req := httptest.NewRequest("GET", "/test", nil)

				ctx := context.WithValue(req.Context(), jwtmiddleware.ContextKey{}, "not-a-validated-claims")
				req = req.WithContext(ctx)

				recorder := httptest.NewRecorder()
				middlewareHandler.ServeHTTP(recorder, req)

				So(recorder.Code, ShouldEqual, http.StatusBadRequest)
				So(recorder.Header().Get("Content-Type"), ShouldEqual, "application/json")
				So(recorder.Body.String(), ShouldEqual, `{"message":"Cannot access token."}`)
			})

			Convey("Test with ValidatedClaims but wrong CustomClaims type", func() {
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

			Convey("Test scope validation failure when missing the required scope", func() {
				req := httptest.NewRequest("GET", "/test", nil)

				claims := &jwtvalidator.ValidatedClaims{
					CustomClaims: &middleware.CustomAuth0Claims{
						Scope: "read:basic write:basic",
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

			Convey("Test successful scope validation", func() {
				req := httptest.NewRequest("GET", "/test", nil)

				claims := &jwtvalidator.ValidatedClaims{
					CustomClaims: &middleware.CustomAuth0Claims{
						Scope: "admin:read write:basic read:basic",
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

		Convey("Should validate Claims interface implementation", func() {
			claims := &middleware.CustomAuth0Claims{
				Scope: "test:scope",
			}

			// Test the Validate method specifically
			err := claims.Validate(context.Background())
			So(err, ShouldBeNil)

			err = claims.Validate(context.TODO())
			So(err, ShouldBeNil)

			// Test with cancelled context
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			err = claims.Validate(ctx)
			So(err, ShouldBeNil)
		})

		Convey("Should handle factory interface correctly", func() {
			// Test that factory implements the interface correctly
			So(factory, ShouldImplement, (*middleware.IAuth0ValidatorFactory)(nil))

			// Test multiple creations from same factory
			config1 := subrouter_config.AuthorizationConfig{
				RequiredScopes: []string{"scope1"},
			}
			config2 := subrouter_config.AuthorizationConfig{
				RequiredScopes: []string{"scope2"},
			}

			validator1 := factory.NewAuth0ScopeValidator(config1)
			validator2 := factory.NewAuth0ScopeValidator(config2)
			So(validator1, ShouldNotBeNil)
			So(validator2, ShouldNotBeNil)
			So(validator1, ShouldNotEqual, validator2)
		})
	})
}
