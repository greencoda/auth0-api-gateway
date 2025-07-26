package auth0_test

import (
	"testing"

	auth0_config "github.com/greencoda/auth0-api-gateway/internal/config/auth0"
	"github.com/greencoda/confiq"
	yaml_loader "github.com/greencoda/confiq/loaders/yaml"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewConfig(t *testing.T) {
	Convey("When creating a new Auth0 config", t, func() {
		Convey("With valid config set", func() {
			var (
				configSet      = confiq.New()
				expectedConfig = auth0_config.Config{
					Audience: "https://test-api.example.com",
					Domain:   "test-tenant.auth0.com",
				}
			)

			err := configSet.Load(
				yaml_loader.Load().FromFile("testdata/valid_config.yaml"),
			)
			So(err, ShouldBeNil)

			config, err := auth0_config.NewConfig(configSet)
			So(err, ShouldBeNil)
			So(config, ShouldNotBeNil)
			So(*config, ShouldResemble, expectedConfig)
		})

		Convey("With empty config, using default values", func() {
			var (
				configSet      = confiq.New()
				expectedConfig = auth0_config.Config{
					Audience: "https://your-auth0-api.yourdomain.io",
					Domain:   "your-auth0-tenant.eu.auth0.com",
				}
			)

			err := configSet.Load(
				yaml_loader.Load().FromFile("testdata/empty_config.yaml"),
			)
			So(err, ShouldBeNil)

			config, err := auth0_config.NewConfig(configSet)
			So(err, ShouldBeNil)
			So(config, ShouldNotBeNil)
			So(*config, ShouldResemble, expectedConfig)
		})

		Convey("With nil config set", func() {
			config, err := auth0_config.NewConfig(nil)
			So(err, ShouldNotBeNil)
			So(config, ShouldBeNil)
		})

		Convey("With empty config set", func() {
			var (
				configSet      = confiq.New()
				expectedConfig = auth0_config.Config{
					Audience: "https://your-auth0-api.yourdomain.io",
					Domain:   "your-auth0-tenant.eu.auth0.com",
				}
			)

			config, err := auth0_config.NewConfig(configSet)
			So(err, ShouldBeNil)
			So(config, ShouldNotBeNil)
			So(*config, ShouldResemble, expectedConfig)
		})
	})
}
