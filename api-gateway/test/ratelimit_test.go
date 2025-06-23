package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"api-gateway/middlewares"
)

func TestRateLimitMiddleware(t *testing.T) {
	h := middlewares.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 12; i++ {
		req := httptest.NewRequest("GET", "/user/profile", nil)
		rw := httptest.NewRecorder()
		h.ServeHTTP(rw, req)
		if i < 10 && rw.Code != http.StatusOK {
			t.Errorf("expected 200 OK, got %d on req %d", rw.Code, i)
		}
		if i >= 10 && rw.Code != http.StatusTooManyRequests {
			t.Errorf("expected 429 Too Many Requests, got %d on req %d", rw.Code, i)
		}
	}
}
