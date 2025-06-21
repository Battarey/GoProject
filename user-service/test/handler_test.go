package test

import (
	"context"
	"testing"
	"user-service/handler"
	"user-service/model"
	user "user-service/proto"
	"user-service/repository"
	"user-service/security"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupHandlerTest() *handler.UserServer {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.User{})
	repo := repository.NewUserRepository(db)
	jwt := security.NewJWTService("testsecret")
	return &handler.UserServer{Repo: repo, JwtService: jwt}
}

func TestRegisterAndLogin(t *testing.T) {
	h := setupHandlerTest()
	ctx := context.Background()

	regReq := &user.RegisterRequest{
		Username: "testuser",
		Email:    "test2@example.com",
		Password: "password123",
	}
	regResp, err := h.Register(ctx, regReq)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if regResp.UserId == "" {
		t.Error("expected non-empty UserId")
	}

	loginReq := &user.LoginRequest{
		Email:    "test2@example.com",
		Password: "password123",
	}
	loginResp, err := h.Login(ctx, loginReq)
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if loginResp.Token == "" {
		t.Error("expected non-empty token")
	}
}
