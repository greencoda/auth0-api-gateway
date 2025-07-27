package cors_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	subrouter_config "github.com/greencoda/auth0-api-gateway/internal/config/subrouter"
	middleware "github.com/greencoda/auth0-api-gateway/internal/middleware/cors"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_CORS_Handler(t *testing.T) {
	Convey("When using CORS handler", t, func() {
		factory := middleware.NewCORSFactory()

		Convey("Should handle preflight requests", func() {
			config := subrouter_config.CORSConfig{
				AllowedOrigins: []string{"https://example.com"},
				AllowedMethods: []string{"GET", "POST"},
				AllowedHeaders: []string{"Authorization"},
				MaxAge:         3600,
			}

			cors := factory.NewCORS(config)
			handler := cors.Handler()

			testHandler := http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
				responseWriter.WriteHeader(http.StatusOK)
				_, _ = responseWriter.Write([]byte("test response"))
			})

			wrappedHandler := handler(testHandler)

			req := httptest.NewRequest("OPTIONS", "http://example.com/api", nil)
			req.Header.Set("Origin", "https://example.com")
			req.Header.Set("Access-Control-Request-Method", "POST")
			req.Header.Set("Access-Control-Request-Headers", "Authorization")
			responseRecorder := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(responseRecorder, req)

			So(responseRecorder.Header().Get("Access-Control-Allow-Origin"), ShouldNotBeEmpty)
		})

		Convey("Should handle actual requests", func() {
			config := subrouter_config.CORSConfig{
				AllowedOrigins: []string{"https://example.com"},
				AllowedMethods: []string{"GET", "POST"},
				AllowedHeaders: []string{"Authorization"},
			}

			cors := factory.NewCORS(config)
			handler := cors.Handler()

			testHandler := http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
				responseWriter.WriteHeader(http.StatusOK)
				_, _ = responseWriter.Write([]byte("success"))
			})

			wrappedHandler := handler(testHandler)

			req := httptest.NewRequest("GET", "http://example.com/api", nil)
			req.Header.Set("Origin", "https://example.com")
			responseRecorder := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(responseRecorder, req)

			So(responseRecorder.Code, ShouldEqual, http.StatusOK)
			So(responseRecorder.Body.String(), ShouldEqual, "success")
		})

		Convey("Should work with wildcard origins", func() {
			config := subrouter_config.CORSConfig{
				AllowedOrigins: []string{"*"},
				AllowedMethods: []string{"GET"},
				AllowedHeaders: []string{"*"},
			}

			cors := factory.NewCORS(config)
			handler := cors.Handler()

			testHandler := http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
				responseWriter.WriteHeader(http.StatusOK)
			})

			wrappedHandler := handler(testHandler)

			req := httptest.NewRequest("GET", "http://example.com/api", nil)
			req.Header.Set("Origin", "https://anydomain.com")
			responseRecorder := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(responseRecorder, req)

			So(responseRecorder.Code, ShouldEqual, http.StatusOK)
		})

		Convey("Should handle multiple allowed origins", func() {
			config := subrouter_config.CORSConfig{
				AllowedOrigins: []string{"https://app1.example.com", "https://app2.example.com"},
				AllowedMethods: []string{"GET", "POST"},
			}

			cors := factory.NewCORS(config)
			handler := cors.Handler()

			testHandler := http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
				responseWriter.WriteHeader(http.StatusOK)
			})

			wrappedHandler := handler(testHandler)

			req1 := httptest.NewRequest("GET", "http://example.com/api", nil)
			req1.Header.Set("Origin", "https://app1.example.com")
			w1 := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(w1, req1)
			So(w1.Code, ShouldEqual, http.StatusOK)

			req2 := httptest.NewRequest("GET", "http://example.com/api", nil)
			req2.Header.Set("Origin", "https://app2.example.com")
			w2 := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(w2, req2)
			So(w2.Code, ShouldEqual, http.StatusOK)
		})

		Convey("Should handle requests without Origin header", func() {
			config := subrouter_config.CORSConfig{
				AllowedOrigins: []string{"https://example.com"},
				AllowedMethods: []string{"GET"},
			}

			cors := factory.NewCORS(config)
			handler := cors.Handler()

			testHandler := http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
				responseWriter.WriteHeader(http.StatusOK)
				_, _ = responseWriter.Write([]byte("no cors needed"))
			})

			wrappedHandler := handler(testHandler)

			req := httptest.NewRequest("GET", "http://example.com/api", nil)
			responseRecorder := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(responseRecorder, req)

			So(responseRecorder.Code, ShouldEqual, http.StatusOK)
			So(responseRecorder.Body.String(), ShouldEqual, "no cors needed")
		})
	})
}
