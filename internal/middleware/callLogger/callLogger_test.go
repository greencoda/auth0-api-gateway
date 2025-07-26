package callLogger_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/callLogger"
	"github.com/rs/zerolog"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewCallLogger(t *testing.T) {
	Convey("When creating a new call logger", t, func() {
		var buf bytes.Buffer
		logger := zerolog.New(&buf)

		Convey("With valid logger", func() {
			params := middleware.CallLoggerParams{
				Logger: logger,
			}

			callLogger := middleware.NewCallLogger(params)
			So(callLogger, ShouldNotBeNil)
			So(callLogger, ShouldImplement, (*middleware.ICallLogger)(nil))
		})
	})
}

func Test_CallLogger_Handler(t *testing.T) {
	Convey("When using call logger handler", t, func() {
		var buf bytes.Buffer
		logger := zerolog.New(&buf)

		params := middleware.CallLoggerParams{
			Logger: logger,
		}

		callLogger := middleware.NewCallLogger(params)

		Convey("Should log requests", func() {
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("test response"))
			})

			wrappedHandler := callLogger.Handler(testHandler)

			req := httptest.NewRequest("GET", "http://example.com/test/path", nil)
			req.RemoteAddr = "192.168.1.1:12345"
			w := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Body.String(), ShouldEqual, "test response")

			logOutput := buf.String()
			So(logOutput, ShouldContainSubstring, "/test/path")
			So(logOutput, ShouldContainSubstring, "192.168.1.1:12345")
		})

		Convey("Should log different paths", func() {
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			wrappedHandler := callLogger.Handler(testHandler)

			testCases := []string{
				"/api/v1/users",
				"/health",
				"/metrics",
				"/",
			}

			for _, path := range testCases {
				Convey("Should log path: "+path, func() {
					buf.Reset()

					req := httptest.NewRequest("GET", "http://example.com"+path, nil)
					req.RemoteAddr = "127.0.0.1:8080"
					w := httptest.NewRecorder()

					wrappedHandler.ServeHTTP(w, req)
					So(w.Code, ShouldEqual, http.StatusOK)
					So(w.Body.String(), ShouldEqual, "")

					logOutput := buf.String()
					So(logOutput, ShouldContainSubstring, path)
					So(logOutput, ShouldContainSubstring, "127.0.0.1:8080")
				})
			}
		})

		Convey("Should handle different HTTP methods", func() {
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			wrappedHandler := callLogger.Handler(testHandler)

			testCases := []string{
				"GET",
				"POST",
				"PUT",
				"DELETE",
				"PATCH",
			}

			for _, method := range testCases {
				Convey("Should handle method: "+method, func() {
					buf.Reset()

					req := httptest.NewRequest(method, "http://example.com/test", nil)
					req.RemoteAddr = "10.0.0.1:9000"
					w := httptest.NewRecorder()

					wrappedHandler.ServeHTTP(w, req)
					So(w.Code, ShouldEqual, http.StatusOK)
					So(w.Body.String(), ShouldEqual, "")

					logOutput := buf.String()
					So(logOutput, ShouldContainSubstring, "/test")
					So(logOutput, ShouldContainSubstring, "10.0.0.1:9000")
				})
			}
		})

		Convey("Should handle requests with query parameters", func() {
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			wrappedHandler := callLogger.Handler(testHandler)

			req := httptest.NewRequest("GET", "http://example.com/api?param=value&other=123", nil)
			req.RemoteAddr = "192.168.0.100:5555"
			w := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Body.String(), ShouldEqual, "")

			logOutput := buf.String()
			So(logOutput, ShouldContainSubstring, "/api")
			So(logOutput, ShouldContainSubstring, "192.168.0.100:5555")
		})

		Convey("Should handle empty RemoteAddr", func() {
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			wrappedHandler := callLogger.Handler(testHandler)

			req := httptest.NewRequest("GET", "http://example.com/test", nil)
			req.RemoteAddr = ""
			w := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Body.String(), ShouldEqual, "")

			logOutput := buf.String()
			So(logOutput, ShouldContainSubstring, "/test")
		})
	})
}
