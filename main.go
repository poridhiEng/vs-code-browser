package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {

	fmt.Println("Running..")

	fqdn := "${splitPath}.${splitPath}.svc.cluster.local"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Count(r.URL.Path, "/") < 2 {
			http.NotFound(w, r)
			return
		}

		splitPaths := strings.SplitN(r.URL.Path, "/", 3)
		fmt.Println("splitting----->", splitPaths)

		splitPath := splitPaths[1] // Extract the desired part of the path
		fmt.Println("splitPath---->", splitPath)

		// Construct the FQDN by replacing placeholders with the splitPath value
		target := strings.ReplaceAll(fqdn, "${splitPath}", splitPath)
		fmt.Println("target---->", target)

		remote, err := url.Parse(target)
		if err != nil {
			panic(err)
		}

		fmt.Println("remote------>", remote)

		proxy := httputil.NewSingleHostReverseProxy(remote)

		// Modify the request before it gets forwarded to the target
		director := proxy.Director
		proxy.Director = func(req *http.Request) {
			director(req)

			log.Println(req.URL.Path)

			// Split the URL path and remove the first element (dynamic text/ID)
			splitPath := strings.SplitN(req.URL.Path, "/", 3)
			log.Println("namespace", splitPath)
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
