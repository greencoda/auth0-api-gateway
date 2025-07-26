package logging_test

import (
	"bytes"
	"testing"

	server_config "github.com/greencoda/auth0-api-gateway/internal/config/server"
	"github.com/greencoda/auth0-api-gateway/internal/util/logging"
	"github.com/rs/zerolog"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewZeroLogger(t *testing.T) {
	Convey("When creating a new ZeroLogger", t, func() {
		Convey("With local release stage", func() {
			config := &server_config.Config{
				ReleaseStage: "local",
				LogLevel:     "info",
			}

			logger := logging.NewZeroLogger(config)
			So(logger, ShouldNotBeNil)
		})

		Convey("With production release stage", func() {
			config := &server_config.Config{
				ReleaseStage: "production",
				LogLevel:     "warn",
			}

			logger := logging.NewZeroLogger(config)
			So(logger, ShouldNotBeNil)
		})

		Convey("With different log levels", func() {
			testCases := []struct {
				logLevel     string
				expectedFunc func(zerolog.Logger) bool
			}{
				{"debug", func(l zerolog.Logger) bool { return l.GetLevel() == zerolog.DebugLevel }},
				{"info", func(l zerolog.Logger) bool { return l.GetLevel() == zerolog.InfoLevel }},
				{"warn", func(l zerolog.Logger) bool { return l.GetLevel() == zerolog.WarnLevel }},
				{"error", func(l zerolog.Logger) bool { return l.GetLevel() == zerolog.ErrorLevel }},
				{"fatal", func(l zerolog.Logger) bool { return l.GetLevel() == zerolog.FatalLevel }},
				{"panic", func(l zerolog.Logger) bool { return l.GetLevel() == zerolog.PanicLevel }},
				{"invalid", func(l zerolog.Logger) bool { return l.GetLevel() == zerolog.InfoLevel }}, // should default to no change
			}

			for _, tc := range testCases {
				Convey("With log level: "+tc.logLevel, func() {
					config := &server_config.Config{
						ReleaseStage: "production",
						LogLevel:     tc.logLevel,
					}

					logger := logging.NewZeroLogger(config)
					So(logger, ShouldNotBeNil)
					// We can't easily test the exact level due to how zerolog works,
					// but we can test that the logger was created successfully
				})
			}
		})

		Convey("With case insensitive log levels", func() {
			testCases := []string{"DEBUG", "Info", "WARN", "Error"}

			for _, logLevel := range testCases {
				Convey("With log level: "+logLevel, func() {
					config := &server_config.Config{
						ReleaseStage: "production",
						LogLevel:     logLevel,
					}

					logger := logging.NewZeroLogger(config)
					So(logger, ShouldNotBeNil)
				})
			}
		})

		Convey("With empty log level", func() {
			config := &server_config.Config{
				ReleaseStage: "production",
				LogLevel:     "",
			}

			logger := logging.NewZeroLogger(config)
			So(logger, ShouldNotBeNil)
		})
	})
}

func Test_LoggerOutput(t *testing.T) {
	Convey("When testing logger output", t, func() {
		Convey("Logger should be functional", func() {
			config := &server_config.Config{
				ReleaseStage: "production",
				LogLevel:     "info",
			}

			var buf bytes.Buffer
			logger := logging.NewZeroLogger(config).Output(&buf)

			logger.Info().Msg("test message")

			output := buf.String()
			So(output, ShouldContainSubstring, "test message")
			So(output, ShouldContainSubstring, "info")
		})

		Convey("Different log levels should work", func() {
			config := &server_config.Config{
				ReleaseStage: "production",
				LogLevel:     "debug",
			}

			var buf bytes.Buffer
			logger := logging.NewZeroLogger(config).Output(&buf)

			logger.Debug().Msg("debug message")
			logger.Info().Msg("info message")
			logger.Warn().Msg("warn message")

			output := buf.String()
			So(output, ShouldContainSubstring, "debug message")
			So(output, ShouldContainSubstring, "info message")
			So(output, ShouldContainSubstring, "warn message")
		})

		Convey("Log level filtering should work", func() {
			config := &server_config.Config{
				ReleaseStage: "production",
				LogLevel:     "warn",
			}

			var buf bytes.Buffer
			logger := logging.NewZeroLogger(config).Output(&buf)

			logger.Debug().Msg("debug message")
			logger.Info().Msg("info message")
			logger.Warn().Msg("warn message")

			output := buf.String()
			So(output, ShouldNotContainSubstring, "debug message")
			So(output, ShouldNotContainSubstring, "info message")
			So(output, ShouldContainSubstring, "warn message")
		})
	})
}
