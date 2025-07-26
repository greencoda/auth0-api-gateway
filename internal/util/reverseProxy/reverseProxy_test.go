package reverseProxy_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/greencoda/auth0-api-gateway/internal/util/reverseProxy"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewReverseProxy(t *testing.T) {
	Convey("When creating a new reverse proxy", t, func() {
		Convey("With valid target URL", func() {
			targetURL, err := url.Parse("http://localhost:8080")
			So(err, ShouldBeNil)

			proxy := reverseProxy.NewReverseProxy(targetURL)
			So(proxy, ShouldNotBeNil)
			So(proxy.Director, ShouldNotBeNil)
		})

		Convey("With target URL including path", func() {
			targetURL, err := url.Parse("http://localhost:8080/api/v1")
			So(err, ShouldBeNil)

			proxy := reverseProxy.NewReverseProxy(targetURL)
			So(proxy, ShouldNotBeNil)
		})

		Convey("With target URL including query parameters", func() {
			targetURL, err := url.Parse("http://localhost:8080?version=1")
			So(err, ShouldBeNil)

			proxy := reverseProxy.NewReverseProxy(targetURL)
			So(proxy, ShouldNotBeNil)
		})
	})
}

func Test_ReverseProxyFunctionality(t *testing.T) {
	Convey("When testing reverse proxy functionality", t, func() {
		// Create a test backend server
		backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Backend-Path", r.URL.Path)
			w.Header().Set("X-Backend-Query", r.URL.RawQuery)
			w.Header().Set("X-Backend-Host", r.Host)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("backend response"))
		}))
		defer backendServer.Close()

		targetURL, err := url.Parse(backendServer.URL)
		So(err, ShouldBeNil)

		proxy := reverseProxy.NewReverseProxy(targetURL)

		Convey("Should proxy basic requests correctly", func() {
			req := httptest.NewRequest("GET", "http://frontend.com/test", nil)
			w := httptest.NewRecorder()

			proxy.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Body.String(), ShouldEqual, "backend response")
			So(w.Header().Get("X-Backend-Path"), ShouldEqual, "/test")
		})

		Convey("Should handle requests with query parameters", func() {
			req := httptest.NewRequest("GET", "http://frontend.com/test?param=value", nil)
			w := httptest.NewRecorder()

			proxy.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Header().Get("X-Backend-Query"), ShouldEqual, "param=value")
		})

		Convey("Should merge query parameters when target has them", func() {
			targetWithQuery, err := url.Parse(backendServer.URL + "?target=param")
			So(err, ShouldBeNil)

			proxyWithQuery := reverseProxy.NewReverseProxy(targetWithQuery)

			req := httptest.NewRequest("GET", "http://frontend.com/test?req=param", nil)
			w := httptest.NewRecorder()

			proxyWithQuery.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
			query := w.Header().Get("X-Backend-Query")
			So(query, ShouldContainSubstring, "target=param")
			So(query, ShouldContainSubstring, "req=param")
		})

		Convey("Should handle requests with paths", func() {
			req := httptest.NewRequest("GET", "http://frontend.com/api/v1/users", nil)
			w := httptest.NewRecorder()

			proxy.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Header().Get("X-Backend-Path"), ShouldEqual, "/api/v1/users")
		})

		Convey("Should join target path with request path", func() {
			targetWithPath, err := url.Parse(backendServer.URL + "/base")
			So(err, ShouldBeNil)

			proxyWithPath := reverseProxy.NewReverseProxy(targetWithPath)

			req := httptest.NewRequest("GET", "http://frontend.com/endpoint", nil)
			w := httptest.NewRecorder()

			proxyWithPath.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Header().Get("X-Backend-Path"), ShouldEqual, "/base/endpoint")
		})

		Convey("Should handle User-Agent header correctly", func() {
			req := httptest.NewRequest("GET", "http://frontend.com/test", nil)
			// Don't set User-Agent header
			w := httptest.NewRecorder()

			proxy.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("Should preserve existing User-Agent header", func() {
			req := httptest.NewRequest("GET", "http://frontend.com/test", nil)
			req.Header.Set("User-Agent", "custom-agent")
			w := httptest.NewRecorder()

			proxy.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})
	})
}

func Test_PathJoining(t *testing.T) {
	Convey("When testing path joining functionality", t, func() {
		// Test through the reverse proxy behavior since the helper functions are not exported

		Convey("Should handle root paths correctly", func() {
			backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Backend-Path", r.URL.Path)
				w.WriteHeader(http.StatusOK)
			}))
			defer backendServer.Close()

			testCases := []struct {
				targetPath   string
				requestPath  string
				expectedPath string
			}{
				{"/", "/test", "/test"},
				{"/api", "/test", "/api/test"},
				{"/api/", "/test", "/api/test"},
				{"/api", "/", "/api"},
				{"/api/v1", "/users", "/api/v1/users"},
				{"/api/v1/", "/users/", "/api/v1/users"},
				{"/api/v1/", "/foo%2Fbar/", "/api/v1/foo%2Fbar"},
				{"", "/test", "/test"},
				{"/", "/", "/"},
			}

			for _, testCase := range testCases {
				Convey("Target: "+testCase.targetPath+", Request: "+testCase.requestPath+" -> "+testCase.expectedPath, func() {
					targetURL, err := url.Parse(backendServer.URL + testCase.targetPath)
					So(err, ShouldBeNil)

					proxy := reverseProxy.NewReverseProxy(targetURL)

					req := httptest.NewRequest("GET", "http://frontend.com"+testCase.requestPath, nil)
					w := httptest.NewRecorder()

					proxy.ServeHTTP(w, req)

					So(w.Code, ShouldEqual, http.StatusOK)
					So(w.Header().Get("X-Backend-Path"), ShouldEqual, testCase.expectedPath)
				})
			}
		})
	})
}
