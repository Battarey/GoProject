package test

import (
	"context"
	"testing"
	"time"

	"task-service/handler"
	"task-service/model"
	"task-service/repository"
	"task-service/security"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestServer создаёт тестовый gRPC сервер с in-memory SQLite
func setupTestServer(t *testing.T) *handler.TaskServer {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.Task{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	repo := repository.NewTaskRepository(db)
	jwtService := security.NewJWTService("testsecret")
	rateLimiter := handler.NewRateLimiter(10 * time.Millisecond)
	return &handler.TaskServer{Repo: repo, JwtService: jwtService, RateLimiter: rateLimiter}
}

// makeJWT генерирует JWT для тестов
func makeJWT(t *testing.T, secret, userID, role string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": userID, "role": role})
	tokStr, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign jwt: %v", err)
	}
	return tokStr
}

// ctxWithJWT возвращает context с JWT в metadata
func ctxWithJWT(token string) context.Context {
	md := metadata.New(map[string]string{"authorization": "Bearer " + token})
	return metadata.NewIncomingContext(context.Background(), md)
}
