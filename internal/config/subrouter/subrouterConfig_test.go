package subrouter_test

import (
	"testing"

	subrouter_config "github.com/greencoda/auth0-api-gateway/internal/config/subrouter"
	"github.com/greencoda/confiq"
	yaml_loader "github.com/greencoda/confiq/loaders/yaml"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewConfig(t *testing.T) {
	Convey("When creating a new Subrouter config", t, func() {
		Convey("With valid config set", func() {
			var (
				configSet      = confiq.New()
				expectedConfig = subrouter_config.Config{
					{
						Name:        "Test API",
						TargetURL:   "http://localhost:8088",
						Prefix:      "/api/v1",
						StripPrefix: true,
						GZip:        true,
					},
					{
						Name:        "Another API",
						TargetURL:   "http://localhost:9090",
						Prefix:      "/api/v2",
						StripPrefix: false,
						GZip:        false,
					},
				}
			)

			err := configSet.Load(
				yaml_loader.Load().FromFile("testdata/valid_config.yaml"),
			)
			So(err, ShouldBeNil)

			config, err := subrouter_config.NewConfig(configSet)
			So(err, ShouldBeNil)
			So(config, ShouldNotBeNil)
			So(*config, ShouldHaveLength, 2)
			So(*config, ShouldResemble, expectedConfig)
		})

		Convey("With default values", func() {
			var (
				configSet      = confiq.New()
				expectedConfig = subrouter_config.Config{
					{
						Name:        "Minimal API",
						TargetURL:   "http://localhost:8080",
						Prefix:      "/minimal",
						StripPrefix: false,
						GZip:        false,
					},
				}
			)

			err := configSet.Load(
				yaml_loader.Load().FromFile("testdata/minimal_config.yaml"),
			)
			So(err, ShouldBeNil)

			config, err := subrouter_config.NewConfig(configSet)
			So(err, ShouldBeNil)
			So(config, ShouldNotBeNil)
			So(*config, ShouldHaveLength, 1)
			So(*config, ShouldResemble, expectedConfig)
		})

		Convey("With nil config set", func() {
			config, err := subrouter_config.NewConfig(nil)
			So(err, ShouldNotBeNil)
			So(config, ShouldBeNil)
		})

		Convey("With empty config set", func() {
			configSet := confiq.New()
			config, err := subrouter_config.NewConfig(configSet)
			So(err, ShouldNotBeNil)
			So(config, ShouldBeNil)
		})

		Convey("With empty subrouters array", func() {
			configSet := confiq.New()
			err := configSet.Load(
				yaml_loader.Load().FromFile("testdata/empty_config.yaml"),
			)
			So(err, ShouldBeNil)

			config, err := subrouter_config.NewConfig(configSet)
			So(err, ShouldNotBeNil)
			So(config, ShouldBeNil)
		})
	})
}
