package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	// Define the target server to proxy to
	target := "http://localhost:8080"
	remote, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)

	// Custom Director to modify the request before forwarding
	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req) // Call the default director

		// Rewrite the path
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/namespace")
		if req.URL.Path == "" {
			req.URL.Path = "/"
		}

		// Set headers
		req.Header.Set("Host", req.Host)
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Header.Set("X-Forwarded-For", req.RemoteAddr)
		req.Header.Set("X-Forwarded-Proto", "http")
		req.Header.Set("X-Real-IP", req.RemoteAddr)
	}

	// Setup the server
	http.HandleFunc("/namespace/", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Host = remote.Host
		r.URL.Scheme = remote.Scheme
		proxy.ServeHTTP(w, r)
	})

	// Start the server
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
