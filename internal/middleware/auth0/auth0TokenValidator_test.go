package auth0_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_ "embed"

	auth0_config "github.com/greencoda/auth0-api-gateway/internal/config/auth0"
	middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/auth0"
	"github.com/h2non/gock"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	//go:embed testdata/mock_jwks.json
	jwks string

	//go:embed testdata/mock_openIdConfig.json
	openIdConfig string

	//go:embed testdata/jwtToken_valid.txt
	validJWTToken string

	//go:embed testdata/jwtToken_invalid.txt
	invalidJWTToken string
)

func Test_Auth0TokenValidator_Interface(t *testing.T) {
	Convey("When testing Auth0 token validator interface", t, func() {
		factory := middleware.NewAuth0ValidatorFactory()

		Convey("Should implement IAuth0TokenValidator interface", func() {
			defer gock.Off()

			gock.New("https://test-auth0.local").
				Get("/.well-known/openid-configuration").
				Reply(200).
				JSON(openIdConfig)

			gock.New("https://test-auth0.local").
				Get("/.well-known/jwks.json").
				Reply(200).
				JSON(jwks)

			config := auth0_config.Config{
				Audience: "https://test-api.local/",
				Domain:   "test-auth0.local",
			}

			validator, err := factory.NewAuth0TokenValidator(config)
			So(err, ShouldBeNil)
			So(validator, ShouldNotBeNil)
			So(validator, ShouldImplement, (*middleware.IAuth0TokenValidator)(nil))

			handler := validator.Handler()
			So(handler, ShouldNotBeNil)

			Convey("Should return a valid middleware function", func() {
				testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte("success"))
				})

				wrappedHandler := handler(testHandler)

				Convey("When a request has a valid JWT", func() {
					req := httptest.NewRequest("GET", "/test", nil)
					req.Header.Set("authorization", "Bearer "+validJWTToken)

					recorder := httptest.NewRecorder()

					wrappedHandler.ServeHTTP(recorder, req)
					So(recorder.Body.String(), ShouldEqual, "success")
					So(recorder.Code, ShouldEqual, http.StatusOK)
				})

				Convey("When a request has an invalid JWT", func() {
					req := httptest.NewRequest("GET", "/test", nil)
					req.Header.Set("Authorization", "Bearer "+invalidJWTToken)

					recorder := httptest.NewRecorder()

					wrappedHandler.ServeHTTP(recorder, req)
					So(recorder.Code, ShouldEqual, http.StatusUnauthorized)
					So(recorder.Header().Get("Content-Type"), ShouldEqual, "application/json")
				})
			})
		})

		Convey("Should handle invalid configuration", func() {
			Convey("With invalid domain URL", func() {
				config := auth0_config.Config{
					Audience: "https://api.example.com",
					Domain:   "invalid domain with spaces",
				}

				validator, err := factory.NewAuth0TokenValidator(config)
				So(err, ShouldNotBeNil)
				So(validator, ShouldBeNil)
			})

			Convey("With empty domain", func() {
				config := auth0_config.Config{
					Audience: "https://api.example.com",
					Domain:   "",
				}

				validator, err := factory.NewAuth0TokenValidator(config)
				So(err, ShouldBeNil)
				So(validator, ShouldNotBeNil)
			})

			Convey("With empty audience", func() {
				config := auth0_config.Config{
					Audience: "",
					Domain:   "example.auth0.com",
				}

				validator, err := factory.NewAuth0TokenValidator(config)
				So(err, ShouldNotBeNil)
				So(validator, ShouldBeNil)
			})
		})

		Convey("Should create validator with different configurations", func() {
			testConfigs := []auth0_config.Config{
				{
					Audience: "https://api1.example.com",
					Domain:   "tenant1.auth0.com",
				},
				{
					Audience: "https://api2.example.com",
					Domain:   "tenant2.eu.auth0.com",
				},
				{
					Audience: "test-audience",
					Domain:   "test.auth0.com",
				},
			}

			for i, config := range testConfigs {
				Convey("Should handle config "+string(rune(i+'1')), func() {
					validator, err := factory.NewAuth0TokenValidator(config)
					So(err, ShouldBeNil)
					So(validator, ShouldNotBeNil)
					So(validator, ShouldImplement, (*middleware.IAuth0TokenValidator)(nil))
				})
			}
		})
	})
}

func Test_BuildJWTMiddlewareFunc_EdgeCases(t *testing.T) {
	Convey("When testing JWT middleware function building", t, func() {
		factory := middleware.NewAuth0ValidatorFactory()

		Convey("Should handle special characters in domain", func() {
			config := auth0_config.Config{
				Audience: "https://api.example.com",
				Domain:   "test-tenant.auth0.com",
			}

			validator, err := factory.NewAuth0TokenValidator(config)
			So(err, ShouldBeNil)
			So(validator, ShouldNotBeNil)
		})

		Convey("Should handle domain with subdomains", func() {
			config := auth0_config.Config{
				Audience: "https://api.example.com",
				Domain:   "subdomain.tenant.eu.auth0.com",
			}

			validator, err := factory.NewAuth0TokenValidator(config)
			So(err, ShouldBeNil)
			So(validator, ShouldNotBeNil)
		})

		Convey("Should handle different audience formats", func() {
			testAudiences := []string{
				"https://api.example.com",
				"api.example.com",
				"my-api",
				"urn:example:api",
			}

			for _, audience := range testAudiences {
				config := auth0_config.Config{
					Audience: audience,
					Domain:   "test.auth0.com",
				}

				validator, err := factory.NewAuth0TokenValidator(config)
				So(err, ShouldBeNil)
				So(validator, ShouldNotBeNil)
			}
		})
	})
}
