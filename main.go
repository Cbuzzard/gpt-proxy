package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	API_KEY := os.Getenv("API_KEY")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	remote, err := url.Parse("https://api.openai.com/v1/completions")
	if err != nil {
		panic(err)
	}

	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			r.Host = remote.Host
			w.Header().Set("Authorization", fmt.Sprintf("Bearer %v", API_KEY))
			p.ServeHTTP(w, r)
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	http.HandleFunc("/", handler(proxy))
	err = http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		panic(err)
	}
}
