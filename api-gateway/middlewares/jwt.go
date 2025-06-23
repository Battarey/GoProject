package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"api-gateway/handlers"
)

// JWTMiddleware проверяет JWT-токен (демо-реализация, без подписи)
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/user/") {
			token := ""
			auth := r.Header.Get("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				token = strings.TrimPrefix(auth, "Bearer ")
			}
			if token == "" {
				handlers.WriteJSONError(w, http.StatusUnauthorized, "Unauthorized: no token")
				return
			}
			parts := strings.Split(token, ".")
			if len(parts) != 3 {
				handlers.WriteJSONError(w, http.StatusUnauthorized, "Invalid token format")
				return
			}
			payload, err := base64.RawURLEncoding.DecodeString(parts[1])
			if err != nil {
				handlers.WriteJSONError(w, http.StatusUnauthorized, "Invalid token payload")
				return
			}
			var claims map[string]interface{}
			_ = json.Unmarshal(payload, &claims)
			if exp, ok := claims["exp"].(float64); ok {
				if int64(exp) < time.Now().Unix() {
					handlers.WriteJSONError(w, http.StatusUnauthorized, "Token expired")
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
