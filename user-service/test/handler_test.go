package test

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"user-service/config"
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

func TestRegisterWithRoleAndValidation(t *testing.T) {
	h := setupHandlerTest()
	ctx := context.Background()

	// Валидная регистрация с ролью admin
	regReq := &user.RegisterRequest{
		Username: "adminuser",
		Email:    "admin@example.com",
		Password: "password123",
		Role:     "admin",
	}
	regResp, err := h.Register(ctx, regReq)
	if err != nil {
		t.Fatalf("register with role failed: %v", err)
	}
	if regResp.UserId == "" {
		t.Error("expected non-empty UserId")
	}

	// Ошибка: невалидный email
	_, err = h.Register(ctx, &user.RegisterRequest{
		Username: "baduser",
		Email:    "bademail",
		Password: "password123",
	})
	if err == nil {
		t.Error("expected error for invalid email")
	}

	// Ошибка: короткий пароль
	_, err = h.Register(ctx, &user.RegisterRequest{
		Username: "shortpass",
		Email:    "short@example.com",
		Password: "123",
	})
	if err == nil {
		t.Error("expected error for short password")
	}

	// Ошибка: невалидная роль
	_, err = h.Register(ctx, &user.RegisterRequest{
		Username: "badrole",
		Email:    "badrole@example.com",
		Password: "password123",
		Role:     "superuser",
	})
	if err == nil {
		t.Error("expected error for invalid role")
	}
}

func TestUpdateAndDeleteUser(t *testing.T) {
	h := setupHandlerTest()
	ctx := context.Background()
	regReq := &user.RegisterRequest{
		Username: "updateuser",
		Email:    "update@example.com",
		Password: "password123",
	}
	regResp, err := h.Register(ctx, regReq)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	updReq := &user.UpdateUserRequest{
		UserId:   regResp.UserId,
		Username: "updated",
		Email:    "updated@example.com",
		Role:     "admin",
	}
	updResp, err := h.UpdateUser(ctx, updReq)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updResp.UserId != regResp.UserId {
		t.Error("user id mismatch after update")
	}
	// Проверка GetProfile
	prof, err := h.GetProfile(ctx, &user.GetProfileRequest{UserId: regResp.UserId})
	if err != nil || prof.Username != "updated" || prof.Role != "admin" {
		t.Error("profile not updated correctly")
	}
	// Удаление
	delResp, err := h.DeleteUser(ctx, &user.DeleteUserRequest{UserId: regResp.UserId})
	if err != nil || !delResp.Success {
		t.Error("delete failed")
	}
}

func TestListUsers(t *testing.T) {
	h := setupHandlerTest()
	ctx := context.Background()
	for i := 1; i <= 5; i++ {
		_, err := h.Register(ctx, &user.RegisterRequest{
			Username: "user" + strconv.Itoa(i),
			Email:    "user" + strconv.Itoa(i) + "@example.com",
			Password: "password123",
			Role:     "user",
		})
		if err != nil {
			t.Fatalf("register failed: %v", err)
		}
	}
	resp, err := h.ListUsers(ctx, &user.ListUsersRequest{Page: 1, PageSize: 3})
	if err != nil {
		t.Fatalf("list users failed: %v", err)
	}
	if len(resp.Users) != 3 {
		t.Errorf("expected 3 users, got %d", len(resp.Users))
	}
}

func TestEmailConfirmation(t *testing.T) {
	h := setupHandlerTest()
	ctx := context.Background()

	// Регистрация пользователя
	regReq := &user.RegisterRequest{
		Username: "mailuser",
		Email:    "mail@example.com",
		Password: "password123",
	}
	_, err := h.Register(ctx, regReq)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Получаем пользователя напрямую из репозитория
	repoUser, err := h.Repo.GetUserByEmail("mail@example.com")
	if err != nil || repoUser == nil {
		t.Fatalf("user not found after register")
	}
	if repoUser.IsEmailConfirmed {
		t.Error("email should not be confirmed after register")
	}
	if repoUser.EmailConfirmationToken == "" {
		t.Error("confirmation token should be set")
	}

	// Подтверждение email с правильным токеном
	confResp, err := h.ConfirmEmail(ctx, &user.ConfirmEmailRequest{
		Email: "mail@example.com",
		Token: repoUser.EmailConfirmationToken,
	})
	if err != nil || !confResp.Success {
		t.Errorf("email confirmation failed: %v, %v", err, confResp.Message)
	}

	// Повторное подтверждение (уже подтверждён)
	confResp2, _ := h.ConfirmEmail(ctx, &user.ConfirmEmailRequest{
		Email: "mail@example.com",
		Token: repoUser.EmailConfirmationToken,
	})
	if confResp2.Success {
		t.Error("should not confirm already confirmed email")
	}

	// Подтверждение с невалидным токеном
	confResp3, _ := h.ConfirmEmail(ctx, &user.ConfirmEmailRequest{
		Email: "mail@example.com",
		Token: "badtoken",
	})
	if confResp3.Success {
		t.Error("should not confirm with invalid token")
	}
}

func TestEmailConfirmation_EdgeCases(t *testing.T) {
	h := setupHandlerTest()
	ctx := context.Background()

	// Подмена отправки email на mock
	handler.SendEmailFunc = func(cfg *config.Config, to, token string) error {
		return nil
	}

	// Попытка подтвердить несуществующий email
	resp, err := h.ConfirmEmail(ctx, &user.ConfirmEmailRequest{
		Email: "notfound@example.com",
		Token: "sometoken",
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.Success {
		t.Error("should not confirm non-existent email")
	}

	// Попытка подтвердить email с пустым токеном
	regReq := &user.RegisterRequest{
		Username: "edgeuser",
		Email:    "edge@example.com",
		Password: "password123",
	}
	_, _ = h.Register(ctx, regReq)
	resp2, _ := h.ConfirmEmail(ctx, &user.ConfirmEmailRequest{
		Email: "edge@example.com",
		Token: "",
	})
	if resp2.Success {
		t.Error("should not confirm with empty token")
	}
}

func TestRateLimiting(t *testing.T) {
	h := setupHandlerTest()
	ctx := context.Background()

	// Проверяем лимит на регистрацию (5 в минуту)
	email := "ratelimit@example.com"
	for i := 0; i < 5; i++ {
		_, err := h.Register(ctx, &user.RegisterRequest{
			Username: "user" + strconv.Itoa(i),
			Email:    email,
			Password: "password123",
		})
		if err != nil && !strings.Contains(err.Error(), "already exists") {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	_, err := h.Register(ctx, &user.RegisterRequest{
		Username: "user6",
		Email:    email,
		Password: "password123",
	})
	if err == nil || !strings.Contains(err.Error(), "too many registration attempts") {
		t.Error("expected rate limit error for registration")
	}

	// Проверяем лимит на логин (10 в минуту)
	loginEmail := "ratelogin@example.com"
	_, _ = h.Register(ctx, &user.RegisterRequest{
		Username: "loginuser",
		Email:    loginEmail,
		Password: "password123",
	})
	for i := 0; i < 10; i++ {
		_, _ = h.Login(ctx, &user.LoginRequest{
			Email:    loginEmail,
			Password: "password123",
		})
	}
	_, err = h.Login(ctx, &user.LoginRequest{
		Email:    loginEmail,
		Password: "password123",
	})
	if err == nil || !strings.Contains(err.Error(), "too many login attempts") {
		t.Error("expected rate limit error for login")
	}
}
