package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api-gateway/middlewares"
)

type errorResp struct {
	Error string `json:"error"`
}

func TestJWTMiddleware_NoToken(t *testing.T) {
	h := middlewares.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/user/profile", nil)
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", rw.Code)
	}
	var resp errorResp
	_ = json.Unmarshal(rw.Body.Bytes(), &resp)
	if resp.Error == "" {
		t.Error("expected error message in JSON response")
	}
}

func TestJWTMiddleware_InvalidFormat(t *testing.T) {
	h := middlewares.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/user/profile", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", rw.Code)
	}
	var resp errorResp
	_ = json.Unmarshal(rw.Body.Bytes(), &resp)
	if resp.Error == "" {
		t.Error("expected error message in JSON response")
	}
}

func TestJWTMiddleware_ExpiredToken(t *testing.T) {
	h := middlewares.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// expired token: header.payload.signature (payload: {"exp":0})
	token := "eyJhbGciOiJIUzI1NiJ9.eyJleHAiOjB9.signature"
	req := httptest.NewRequest("GET", "/user/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", rw.Code)
	}
	var resp errorResp
	_ = json.Unmarshal(rw.Body.Bytes(), &resp)
	if resp.Error == "" {
		t.Error("expected error message in JSON response")
	}
}
