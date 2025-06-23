package handlers

import (
	"fmt"
	"net/http"
)

// HealthHandler возвращает 200 OK
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
