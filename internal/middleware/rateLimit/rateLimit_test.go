package rateLimit_test

import (
	"net/http"
	"net/http/httptest"
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

func Test_RateLimit_Handler(t *testing.T) {
	Convey("When using rate limit handler", t, func() {
		factory := middleware.NewRateLimitFactory()

		Convey("Should allow requests within limit", func() {
			config := subrouter_config.RateLimitConfig{
				Limit:  5,
				Period: time.Minute,
			}

			rateLimit := factory.NewRateLimit(config)
			handler := rateLimit.Handler()

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("success"))
			})

			wrappedHandler := handler(testHandler)

			// Make requests within the limit
			for i := 0; i < 5; i++ {
				req := httptest.NewRequest("GET", "http://example.com/api", nil)
				req.RemoteAddr = "192.168.1.1:8080"
				w := httptest.NewRecorder()

				wrappedHandler.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusOK)
				So(w.Body.String(), ShouldEqual, "success")
			}
		})

		Convey("Should handle different IP addresses separately", func() {
			config := subrouter_config.RateLimitConfig{
				Limit:  2,
				Period: time.Minute,
			}

			rateLimit := factory.NewRateLimit(config)
			handler := rateLimit.Handler()

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("success"))
			})

			wrappedHandler := handler(testHandler)

			// Test first IP
			for i := 0; i < 2; i++ {
				req := httptest.NewRequest("GET", "http://example.com/api", nil)
				req.RemoteAddr = "192.168.1.1:8080"
				w := httptest.NewRecorder()

				wrappedHandler.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, http.StatusOK)
			}

			// Test second IP (should also be allowed)
			for i := 0; i < 2; i++ {
				req := httptest.NewRequest("GET", "http://example.com/api", nil)
				req.RemoteAddr = "192.168.1.2:8080"
				w := httptest.NewRecorder()

				wrappedHandler.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, http.StatusOK)
			}
		})

		Convey("Should work with trust forward header", func() {
			config := subrouter_config.RateLimitConfig{
				Limit:              3,
				Period:             time.Minute,
				TrustForwardHeader: true,
			}

			rateLimit := factory.NewRateLimit(config)
			handler := rateLimit.Handler()

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("success"))
			})

			wrappedHandler := handler(testHandler)

			// Make requests with X-Forwarded-For header
			for i := 0; i < 3; i++ {
				req := httptest.NewRequest("GET", "http://example.com/api", nil)
				req.Header.Set("X-Forwarded-For", "10.0.0.1")
				req.RemoteAddr = "192.168.1.1:8080" // This should be ignored when trust forward header is true
				w := httptest.NewRecorder()

				wrappedHandler.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, http.StatusOK)
			}
		})

		Convey("Should handle missing RemoteAddr gracefully", func() {
			config := subrouter_config.RateLimitConfig{
				Limit:  1,
				Period: time.Minute,
			}

			rateLimit := factory.NewRateLimit(config)
			handler := rateLimit.Handler()

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("success"))
			})

			wrappedHandler := handler(testHandler)

			req := httptest.NewRequest("GET", "http://example.com/api", nil)
			req.RemoteAddr = ""
			w := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(w, req)
			// Should handle gracefully, exact behavior depends on underlying limiter
		})

		Convey("Should work with different HTTP methods", func() {
			config := subrouter_config.RateLimitConfig{
				Limit:  3,
				Period: time.Minute,
			}

			rateLimit := factory.NewRateLimit(config)
			handler := rateLimit.Handler()

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(r.Method))
			})

			wrappedHandler := handler(testHandler)

			methods := []string{"GET", "POST", "PUT"}

			for _, method := range methods {
				req := httptest.NewRequest(method, "http://example.com/api", nil)
				req.RemoteAddr = "192.168.1.100:9000"
				w := httptest.NewRecorder()

				wrappedHandler.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, http.StatusOK)
				So(w.Body.String(), ShouldEqual, method)
			}
		})
	})
}
