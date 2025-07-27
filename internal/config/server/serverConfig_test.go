package server_test

import (
	"testing"
	"time"

	server_config "github.com/greencoda/auth0-api-gateway/internal/config/server"
	"github.com/greencoda/confiq"
	yaml_loader "github.com/greencoda/confiq/loaders/yaml"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewConfig(t *testing.T) {
	Convey("When creating a new Server config", t, func() {
		Convey("With valid config set", func() {
			var (
				configSet      = confiq.New()
				expectedConfig = server_config.Config{
					Address:        ":8080",
					ReadTimeout:    30 * time.Second,
					WriteTimeout:   30 * time.Second,
					IdleTimeout:    60 * time.Second,
					MaxHeaderBytes: 2097152,
					ReleaseStage:   "production",
					LogRequests:    true,
					LogLevel:       "debug",
				}
			)

			err := configSet.Load(
				yaml_loader.Load().FromFile("testdata/valid_config.yaml"),
			)
			So(err, ShouldBeNil)

			config, err := server_config.NewConfig(configSet)
			So(err, ShouldBeNil)
			So(config, ShouldNotBeNil)
			So(*config, ShouldResemble, expectedConfig)
		})

		Convey("With empty config, using default values", func() {
			var (
				configSet      = confiq.New()
				expectedConfig = server_config.Config{
					Address:        ":80",
					ReadTimeout:    15 * time.Second,
					WriteTimeout:   15 * time.Second,
					IdleTimeout:    15 * time.Second,
					MaxHeaderBytes: 1048576,
					ReleaseStage:   "local",
					LogRequests:    false,
					LogLevel:       "info",
				}
			)

			err := configSet.Load(
				yaml_loader.Load().FromFile("testdata/empty_config.yaml"),
			)
			So(err, ShouldBeNil)

			config, err := server_config.NewConfig(configSet)
			So(err, ShouldBeNil)
			So(config, ShouldNotBeNil)
			So(*config, ShouldResemble, expectedConfig)
		})

		Convey("With nil config set", func() {
			config, err := server_config.NewConfig(nil)
			So(err, ShouldNotBeNil)
			So(config, ShouldBeNil)
		})

		Convey("With empty config set", func() {
			var (
				configSet      = confiq.New()
				expectedConfig = server_config.Config{
					Address:        ":80",
					ReadTimeout:    15 * time.Second,
					WriteTimeout:   15 * time.Second,
					IdleTimeout:    15 * time.Second,
					MaxHeaderBytes: 1048576,
					ReleaseStage:   "local",
					LogRequests:    false,
					LogLevel:       "info",
				}
			)

			config, err := server_config.NewConfig(configSet)
			So(err, ShouldBeNil)
			So(config, ShouldNotBeNil)
			So(*config, ShouldResemble, expectedConfig)
		})

		Convey("With partial config", func() {
			var (
				configSet      = confiq.New()
				expectedConfig = server_config.Config{
					Address:        ":9000",
					ReadTimeout:    15 * time.Second,
					WriteTimeout:   15 * time.Second,
					IdleTimeout:    15 * time.Second,
					MaxHeaderBytes: 1048576,
					ReleaseStage:   "local",
					LogRequests:    true,
					LogLevel:       "info",
				}
			)

			err := configSet.Load(
				yaml_loader.Load().FromFile("testdata/partial_config.yaml"),
			)
			So(err, ShouldBeNil)

			config, err := server_config.NewConfig(configSet)
			So(err, ShouldBeNil)
			So(config, ShouldNotBeNil)
			So(*config, ShouldResemble, expectedConfig)
		})
	})
}
