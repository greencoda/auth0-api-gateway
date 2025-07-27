package rateLimit_test

import (
	"testing"
	"time"

	subrouter_config "github.com/greencoda/auth0-api-gateway/internal/config/subrouter"
	middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/rateLimit"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewRateLimitFactory(t *testing.T) {
	Convey("When creating a new rate limit factory", t, func() {
		factory := middleware.NewRateLimitFactory()
		So(factory, ShouldNotBeNil)
		So(factory, ShouldImplement, (*middleware.IRateLimitFactory)(nil))
	})
}

func Test_RateLimitFactory_NewRateLimit(t *testing.T) {
	Convey("When creating rate limit middleware", t, func() {
		factory := middleware.NewRateLimitFactory()

		Convey("With basic configuration", func() {
			config := subrouter_config.RateLimitConfig{
				Limit:  10,
				Period: time.Minute,
			}

			rateLimit := factory.NewRateLimit(config)
			So(rateLimit, ShouldNotBeNil)
			So(rateLimit, ShouldImplement, (*middleware.IRateLimit)(nil))

			handler := rateLimit.Handler()
			So(handler, ShouldNotBeNil)
		})

		Convey("With trust forward header enabled", func() {
			config := subrouter_config.RateLimitConfig{
				Limit:              100,
				Period:             time.Second,
				TrustForwardHeader: true,
			}

			rateLimit := factory.NewRateLimit(config)
			So(rateLimit, ShouldNotBeNil)

			handler := rateLimit.Handler()
			So(handler, ShouldNotBeNil)
		})

		Convey("With different time periods", func() {
			testCases := []struct {
				period time.Duration
				name   string
			}{
				{time.Second, "1 second"},
				{time.Minute, "1 minute"},
				{time.Hour, "1 hour"},
				{5 * time.Second, "5 seconds"},
			}

			for _, tc := range testCases {
				Convey("With period: "+tc.name, func() {
					config := subrouter_config.RateLimitConfig{
						Limit:  50,
						Period: tc.period,
					}

					rateLimit := factory.NewRateLimit(config)
					So(rateLimit, ShouldNotBeNil)

					handler := rateLimit.Handler()
					So(handler, ShouldNotBeNil)
				})
			}
		})

		Convey("With different limits", func() {
			testCases := []int64{1, 5, 10, 100, 1000}

			for _, limit := range testCases {
				Convey("With limit: "+string(rune(limit+'0')), func() {
					config := subrouter_config.RateLimitConfig{
						Limit:  limit,
						Period: time.Minute,
					}

					rateLimit := factory.NewRateLimit(config)
					So(rateLimit, ShouldNotBeNil)

					handler := rateLimit.Handler()
					So(handler, ShouldNotBeNil)
				})
			}
		})

		Convey("With zero limit", func() {
			config := subrouter_config.RateLimitConfig{
				Limit:  0,
				Period: time.Minute,
			}

			rateLimit := factory.NewRateLimit(config)
			So(rateLimit, ShouldNotBeNil)

			handler := rateLimit.Handler()
			So(handler, ShouldNotBeNil)
		})
	})
}
