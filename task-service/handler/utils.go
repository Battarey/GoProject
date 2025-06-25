package handler

import (
	"context"
	"strings"
	"task-service/security"

	"google.golang.org/grpc/metadata"
)

// AuthContext содержит user_id и роль, извлечённые из JWT
// Можно расширить по необходимости
func GetAuthContext(ctx context.Context, jwtService *security.JWTService) (userID, role string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", nil
	}
	authHeaders := md["authorization"]
	if len(authHeaders) == 0 {
		return "", "", nil
	}
	token := strings.TrimPrefix(authHeaders[0], "Bearer ")
	claims, err := jwtService.ValidateToken(token)
	if err != nil {
		return "", "", err
	}
	uid, _ := claims["user_id"].(string)
	role, _ = claims["role"].(string)
	return uid, role, nil
}
