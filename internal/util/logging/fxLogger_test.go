package logging_test

import (
	"testing"

	server_config "github.com/greencoda/auth0-api-gateway/internal/config/server"
	"github.com/greencoda/auth0-api-gateway/internal/util/logging"
	"github.com/rs/zerolog"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxevent"
)

func Test_NewFXLogger(t *testing.T) {
	Convey("When creating a new FX logger", t, func() {
		logger := zerolog.New(zerolog.NewConsoleWriter())

		Convey("With production release stage", func() {
			config := &server_config.Config{
				ReleaseStage: "prod",
				LogLevel:     "info",
			}

			fxLogger := logging.NewFXLogger(config, logger)
			So(fxLogger, ShouldNotBeNil)
			So(fxLogger, ShouldEqual, fxevent.NopLogger)
		})

		Convey("With local release stage", func() {
			config := &server_config.Config{
				ReleaseStage: "local",
				LogLevel:     "debug",
			}

			fxLogger := logging.NewFXLogger(config, logger)
			So(fxLogger, ShouldNotBeNil)
			So(fxLogger, ShouldNotEqual, fxevent.NopLogger)
		})

		Convey("With development release stage", func() {
			config := &server_config.Config{
				ReleaseStage: "development",
				LogLevel:     "debug",
			}

			fxLogger := logging.NewFXLogger(config, logger)
			So(fxLogger, ShouldNotBeNil)
			So(fxLogger, ShouldEqual, fxevent.NopLogger)
		})

		Convey("With staging release stage", func() {
			config := &server_config.Config{
				ReleaseStage: "staging",
				LogLevel:     "info",
			}

			fxLogger := logging.NewFXLogger(config, logger)
			So(fxLogger, ShouldNotBeNil)
			So(fxLogger, ShouldEqual, fxevent.NopLogger)
		})

		Convey("With empty release stage", func() {
			config := &server_config.Config{
				ReleaseStage: "",
				LogLevel:     "info",
			}

			fxLogger := logging.NewFXLogger(config, logger)
			So(fxLogger, ShouldNotBeNil)
			So(fxLogger, ShouldEqual, fxevent.NopLogger)
		})
	})
}
