package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	target := "http://localhost:8080"
	remote, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)

	// Modify the request before it gets forwarded to the target
	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)

		// Split the URL path and remove the first element (dynamic text/ID)
		splitPath := strings.SplitN(req.URL.Path, "/", 3)
		if len(splitPath) > 2 {
			req.URL.Path = "/" + splitPath[2]
		} else {
			req.URL.Path = "/"
		}

		// Update the headers
		req.Header.Set("Host", req.Host)
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Forwarded-For", req.RemoteAddr)
		req.Header.Set("X-Forwarded-Proto", req.URL.Scheme)
		req.Header.Set("X-Real-IP", req.RemoteAddr)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Count(r.URL.Path, "/") < 2 {
			http.NotFound(w, r)
			return
		}
		r.URL.Host = remote.Host
		r.URL.Scheme = remote.Scheme
		proxy.ServeHTTP(w, r)
	})

	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
