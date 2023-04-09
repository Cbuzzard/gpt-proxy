package main

import (
	"fmt"
	"io"
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
			if r.Header.Get("Accept") == "text/event-stream" {
				handleSSE(p, w, r)
			} else {
				p.ServeHTTP(w, r)
			}
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	http.HandleFunc("/", handler(proxy))
	err = http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		panic(err)
	}
}

func handleSSE(proxy *httputil.ReverseProxy, w http.ResponseWriter, r *http.Request) {
	proxy.Director(r)
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, "Error contacting backend server", http.StatusInternalServerError)
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Copy headers
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}

	// Set the headers related to SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	w.WriteHeader(resp.StatusCode)

	// Stream the data
	io.Copy(w, resp.Body)
}
