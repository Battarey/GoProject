package test

import (
	"context"
	"testing"
	"user-service/handler"
	"user-service/model"
	pb "user-service/proto"
	"user-service/repository"
	"user-service/security"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRegisterHandler(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.User{})
	repo := repository.NewUserRepository(db)
	jwtService := security.NewJWTService("testsecret")
	h := &handler.UserServer{Repo: repo, JwtService: jwtService}

	resp, err := h.Register(context.Background(), &pb.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "123456",
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.UserId)
}
