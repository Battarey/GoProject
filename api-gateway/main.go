package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	addr := ":8080"
	if v := os.Getenv("GATEWAY_PORT"); v != "" {
		addr = ":" + v
	}
	mux := http.NewServeMux()

	userServiceURL, _ := url.Parse("http://user-service:8081")
	userProxy := httputil.NewSingleHostReverseProxy(userServiceURL)
	mux.Handle("/user/", http.StripPrefix("/user", userProxy))

	log.Printf("api-gateway started on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
