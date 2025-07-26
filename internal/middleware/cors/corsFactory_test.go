package cors_test

import (
	"testing"

	subrouter_config "github.com/greencoda/auth0-api-gateway/internal/config/subrouter"
	middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/cors"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewCORSFactory(t *testing.T) {
	Convey("When creating a new CORS factory", t, func() {
		factory := middleware.NewCORSFactory()
		So(factory, ShouldNotBeNil)
		So(factory, ShouldImplement, (*middleware.ICORSFactory)(nil))
	})
}

func Test_CORSFactory_NewCORS(t *testing.T) {
	Convey("When creating CORS middleware", t, func() {
		factory := middleware.NewCORSFactory()

		Convey("With basic configuration", func() {
			config := subrouter_config.CORSConfig{
				AllowedOrigins: []string{"*"},
				AllowedMethods: []string{"GET", "POST"},
				AllowedHeaders: []string{"Authorization", "Content-Type"},
				MaxAge:         3600,
			}

			cors := factory.NewCORS(config)
			So(cors, ShouldNotBeNil)
			So(cors, ShouldImplement, (*middleware.ICORS)(nil))

			handler := cors.Handler()
			So(handler, ShouldNotBeNil)
		})

		Convey("With credentials allowed", func() {
			config := subrouter_config.CORSConfig{
				AllowedOrigins:   []string{"https://example.com"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
				AllowedHeaders:   []string{"Authorization"},
				AllowCredentials: true,
				MaxAge:           86400,
			}

			cors := factory.NewCORS(config)
			So(cors, ShouldNotBeNil)

			handler := cors.Handler()
			So(handler, ShouldNotBeNil)
		})

		Convey("With exposed headers", func() {
			config := subrouter_config.CORSConfig{
				AllowedOrigins: []string{"https://app.example.com"},
				AllowedMethods: []string{"GET"},
				AllowedHeaders: []string{"Authorization"},
				ExposedHeaders: []string{"X-Custom-Header", "X-Another-Header"},
				MaxAge:         1800,
			}

			cors := factory.NewCORS(config)
			So(cors, ShouldNotBeNil)

			handler := cors.Handler()
			So(handler, ShouldNotBeNil)
		})

		Convey("With empty configuration", func() {
			config := subrouter_config.CORSConfig{}

			cors := factory.NewCORS(config)
			So(cors, ShouldNotBeNil)

			handler := cors.Handler()
			So(handler, ShouldNotBeNil)
		})
	})
}
