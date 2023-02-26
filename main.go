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
	AUTH_KEY := os.Getenv("AUTH_KEY")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	remote, err := url.Parse("https://api.openai.com")
	if err != nil {
		panic(err)
	}

	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Auth") != AUTH_KEY {
				return
			}
			r.Host = remote.Host
			r.Header.Add("Authorization", fmt.Sprintf("Bearer %v", API_KEY))
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
