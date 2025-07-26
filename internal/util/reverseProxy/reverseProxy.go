package reverseProxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

func joinPathsWithSlash(pathA, pathB string) string {
	var pathSegments []string

	for _, path := range []string{pathA, pathB} {
		pathTrimmed := strings.Trim(path, "/")

		if len(pathTrimmed) > 0 {
			pathSegments = append(pathSegments, pathTrimmed)
		}
	}

	return "/" + strings.Join(pathSegments, "/")
}

func joinURLPath(urlA, urlB *url.URL) (path, rawpath string) {
	log.Printf("%s + %s", urlA.RawPath, urlB.RawPath)

	if urlA.RawPath == "" && urlB.RawPath == "" {
		return joinPathsWithSlash(urlA.Path, urlB.Path), ""
	}

	return joinPathsWithSlash(urlA.EscapedPath(), urlB.EscapedPath()), ""
}
