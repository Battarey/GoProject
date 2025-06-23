package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"api-gateway/handlers"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	rw := httptest.NewRecorder()

	handlers.HealthHandler(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rw.Code)
	}
	if rw.Body.String() != "OK" {
		t.Errorf("expected body 'OK', got '%s'", rw.Body.String())
	}
}
