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

func Test_RateLimit_Handler(t *testing.T) {
	Convey("When using rate limit handler", t, func() {
		const remoteAddr = "192.168.1.1:8080"

		factory := middleware.NewRateLimitFactory()

		Convey("Should allow requests within limit", func() {
			config := subrouter_config.RateLimitConfig{
				Limit:  5,
				Period: time.Minute,
			}

			rateLimit := factory.NewRateLimit(config)
			handler := rateLimit.Handler()

			testHandler := http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
				responseWriter.WriteHeader(http.StatusOK)
				_, _ = responseWriter.Write([]byte("success"))
			})

			wrappedHandler := handler(testHandler)

			// Make requests within the limit
			for range 5 {
				req := httptest.NewRequest("GET", "http://example.com/api", nil)
				req.RemoteAddr = remoteAddr
				responseRecorder := httptest.NewRecorder()

				wrappedHandler.ServeHTTP(responseRecorder, req)

				So(responseRecorder.Code, ShouldEqual, http.StatusOK)
				So(responseRecorder.Body.String(), ShouldEqual, "success")
			}
		})

		Convey("Should handle different IP addresses separately", func() {
			config := subrouter_config.RateLimitConfig{
				Limit:  2,
				Period: time.Minute,
			}

			rateLimit := factory.NewRateLimit(config)
			handler := rateLimit.Handler()

			testHandler := http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
				responseWriter.WriteHeader(http.StatusOK)
				_, _ = responseWriter.Write([]byte("success"))
			})

			wrappedHandler := handler(testHandler)

			for range 2 {
				req := httptest.NewRequest("GET", "http://example.com/api", nil)
				req.RemoteAddr = remoteAddr
				responseRecorder := httptest.NewRecorder()

				wrappedHandler.ServeHTTP(responseRecorder, req)
				So(responseRecorder.Code, ShouldEqual, http.StatusOK)
			}

			for range 2 {
				req := httptest.NewRequest("GET", "http://example.com/api", nil)
				req.RemoteAddr = "192.168.1.2:8080"
				responseRecorder := httptest.NewRecorder()

				wrappedHandler.ServeHTTP(responseRecorder, req)
				So(responseRecorder.Code, ShouldEqual, http.StatusOK)
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

			testHandler := http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
				responseWriter.WriteHeader(http.StatusOK)
				_, _ = responseWriter.Write([]byte("success"))
			})

			wrappedHandler := handler(testHandler)

			for range 3 {
				req := httptest.NewRequest("GET", "http://example.com/api", nil)
				req.Header.Set("X-Forwarded-For", "10.0.0.1")
				req.RemoteAddr = remoteAddr
				responseRecorder := httptest.NewRecorder()

				wrappedHandler.ServeHTTP(responseRecorder, req)
				So(responseRecorder.Code, ShouldEqual, http.StatusOK)
			}
		})

		Convey("Should handle missing RemoteAddr gracefully", func() {
			config := subrouter_config.RateLimitConfig{
				Limit:  1,
				Period: time.Minute,
			}

			rateLimit := factory.NewRateLimit(config)
			handler := rateLimit.Handler()

			testHandler := http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
				responseWriter.WriteHeader(http.StatusOK)
				_, _ = responseWriter.Write([]byte("success"))
			})

			wrappedHandler := handler(testHandler)

			req := httptest.NewRequest("GET", "http://example.com/api", nil)
			req.RemoteAddr = ""
			responseRecorder := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(responseRecorder, req)
		})

		Convey("Should work with different HTTP methods", func() {
			config := subrouter_config.RateLimitConfig{
				Limit:  3,
				Period: time.Minute,
			}

			rateLimit := factory.NewRateLimit(config)
			handler := rateLimit.Handler()

			testHandler := http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
				responseWriter.WriteHeader(http.StatusOK)
				_, _ = responseWriter.Write([]byte(req.Method))
			})

			wrappedHandler := handler(testHandler)

			methods := []string{"GET", "POST", "PUT"}

			for _, method := range methods {
				req := httptest.NewRequest(method, "http://example.com/api", nil)
				req.RemoteAddr = "192.168.1.100:9000"
				responseRecorder := httptest.NewRecorder()

				wrappedHandler.ServeHTTP(responseRecorder, req)
				So(responseRecorder.Code, ShouldEqual, http.StatusOK)
				So(responseRecorder.Body.String(), ShouldEqual, method)
			}
		})
	})
}
