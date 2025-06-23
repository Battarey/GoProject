package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// NewUserProxy возвращает reverse proxy для user-service с базовой валидацией
func NewUserProxy() http.Handler {
	userServiceURL, _ := url.Parse("http://user-service:8081")
	userProxy := httputil.NewSingleHostReverseProxy(userServiceURL)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Базовая валидация: только /profile и /register разрешены (пример)
		if !(strings.HasPrefix(r.URL.Path, "/user/profile") || strings.HasPrefix(r.URL.Path, "/user/register")) {
			WriteJSONError(w, http.StatusBadRequest, "invalid user endpoint")
			return
		}
		// Можно добавить доп. валидацию параметров запроса
		// Проксируем дальше
		userProxy.ServeHTTP(w, r)
	})
}
