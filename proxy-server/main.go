package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {

	// this should come from db
	targets := map[string]string{
		"ns1": "http://localhost:8080",
		"ns2": "http://localhost:7080",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Count(r.URL.Path, "/") < 2 {
			http.NotFound(w, r)
			return
		}

		splitPaths := strings.SplitN(r.URL.Path, "/", 3)

		target, _ := targets[splitPaths[1]]

		remote, err := url.Parse(target)
		if err != nil {
			panic(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)

		// Modify the request before it gets forwarded to the target
		director := proxy.Director
		proxy.Director = func(req *http.Request) {
			director(req)

			log.Println(req.URL.Path)

			// Split the URL path and remove the first element (dynamic text/ID)
			splitPath := strings.SplitN(req.URL.Path, "/", 3)
			log.Println(splitPath)
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

		r.URL.Host = remote.Host
		r.URL.Scheme = remote.Scheme
		proxy.ServeHTTP(w, r)
	})

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
