package main

import (
	"log"
	"net/http"
	"os"

	"api-gateway/handlers"
	"api-gateway/middlewares"
)

func main() {
	addr := ":8080"
	if v := os.Getenv("GATEWAY_PORT"); v != "" {
		addr = ":" + v
	}
	mux := http.NewServeMux()

	// Healthcheck endpoint
	mux.HandleFunc("/health", handlers.HealthHandler)

	// /user/* с JWT и rate limiting
	userHandler := handlers.NewUserProxy()
	mux.Handle("/user/", middlewares.JWTMiddleware(middlewares.RateLimitMiddleware(userHandler)))

	// Можно добавить другие сервисы: /task/, /chat/ и т.д.

	// Оборачиваем всё в CORS
	handler := middlewares.CORSMiddleware(mux)

	log.Printf("api-gateway started on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
