package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"api-gateway/middlewares"
)

func TestCORSMiddleware_OPTIONS(t *testing.T) {
	h := middlewares.CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("OPTIONS", "/user/", nil)
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusNoContent {
		t.Errorf("expected 204 No Content, got %d", rw.Code)
	}
	if rw.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS header missing")
	}
}
