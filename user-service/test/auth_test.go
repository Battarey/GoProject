package test

import (
	"context"
	"strconv"
	"testing"
	"time"
	user "user-service/proto"
)

func TestRegisterAndLogin(t *testing.T) {
	h := SetupHandlerTest()
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
	h := SetupHandlerTest()
	ctx := context.Background()

	tests := []struct {
		name      string
		req       *user.RegisterRequest
		expectErr bool
	}{
		{
			name: "valid request",
			req: &user.RegisterRequest{
				Username: "admin",
				Email:    "admin@example.com",
				Password: "password123",
				Role:     "admin",
			},
			expectErr: false,
		},
		{
			name: "missing fields",
			req: &user.RegisterRequest{
				Username: "",
				Email:    "user@example.com",
				Password: "password123",
				Role:     "user",
			},
			expectErr: true,
		},
		{
			name: "invalid email",
			req: &user.RegisterRequest{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "password123",
				Role:     "user",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := h.Register(ctx, tt.req)
			if (err != nil) != tt.expectErr {
				t.Errorf("unexpected error status: got %v, want %v", err != nil, tt.expectErr)
			}
			if !tt.expectErr && resp.UserId == "" {
				t.Error("expected non-empty UserId")
			}
		})
	}
}

func TestEmailConfirmation(t *testing.T) {
	h := SetupHandlerTest()
	ctx := context.Background()

	regReq := &user.RegisterRequest{
		Username: "mailuser",
		Email:    "mail@example.com",
		Password: "password123",
	}
	_, err := h.Register(ctx, regReq)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	repoUser, err := h.Repo.GetUserByEmail("mail@example.com")
	if err != nil || repoUser == nil {
		t.Fatalf("user not found after register")
	}
	if repoUser.EmailConfirmationToken == "" {
		t.Fatal("confirmation token should be set")
	}

	confReq := &user.ConfirmEmailRequest{
		Email: "mail@example.com",
		Token: repoUser.EmailConfirmationToken,
	}
	confResp, err := h.ConfirmEmail(ctx, confReq)
	if err != nil || !confResp.Success {
		t.Errorf("email confirmation failed: %v, %v", err, confResp.Message)
	}
}

func TestEmailConfirmation_EdgeCases(t *testing.T) {
	h := SetupHandlerTest()
	ctx := context.Background()

	confReq := &user.ConfirmEmailRequest{
		Email: "notfound@example.com",
		Token: "sometoken",
	}
	resp, err := h.ConfirmEmail(ctx, confReq)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.Success {
		t.Error("should not confirm non-existent email")
	}

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
	h := SetupHandlerTest()
	ctx := context.Background()

	email := "ratelimit@example.com"
	for i := 0; i < 5; i++ {
		_, err := h.Register(ctx, &user.RegisterRequest{
			Username: "user" + strconv.Itoa(i),
			Email:    email,
			Password: "password123",
		})
		if err != nil && err.Error() != "email already exists" {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	_, err := h.Register(ctx, &user.RegisterRequest{
		Username: "user6",
		Email:    email,
		Password: "password123",
	})
	if err == nil || err.Error() != "too many registration attempts, try later" {
		t.Error("expected rate limit error for registration")
	}
}

func TestPasswordReset(t *testing.T) {
	h := SetupHandlerTest()
	ctx := context.Background()

	email := "reset@example.com"
	_, err := h.Register(ctx, &user.RegisterRequest{
		Username: "resetuser",
		Email:    email,
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	resp, err := h.RequestPasswordReset(ctx, &user.RequestPasswordResetRequest{Email: email})
	if err != nil || !resp.Success {
		t.Fatalf("request password reset failed: %v, %v", err, resp.Message)
	}

	repoUser, _ := h.Repo.GetUserByEmail(email)
	token := repoUser.PasswordResetToken
	if token == "" {
		t.Fatal("reset token not set")
	}

	resetResp, err := h.ResetPassword(ctx, &user.ResetPasswordRequest{
		Email:       email,
		Token:       token,
		NewPassword: "newpassword123",
	})
	if err != nil || !resetResp.Success {
		t.Fatalf("reset password failed: %v, %v", err, resetResp.Message)
	}

	resetResp2, _ := h.ResetPassword(ctx, &user.ResetPasswordRequest{
		Email:       email,
		Token:       token,
		NewPassword: "anotherpass",
	})
	if resetResp2.Success {
		t.Error("should not reset with used token")
	}

	resetResp3, _ := h.ResetPassword(ctx, &user.ResetPasswordRequest{
		Email:       email,
		Token:       "badtoken",
		NewPassword: "anotherpass",
	})
	if resetResp3.Success {
		t.Error("should not reset with invalid token")
	}

	_ = h.Repo.SetPasswordResetToken(email, "shorttoken", time.Now().Add(30*time.Minute).Unix())
	resetResp4, _ := h.ResetPassword(ctx, &user.ResetPasswordRequest{
		Email:       email,
		Token:       "shorttoken",
		NewPassword: "123",
	})
	if resetResp4.Success {
		t.Error("should not reset with short password")
	}

	_ = h.Repo.SetPasswordResetToken(email, "expiredtoken", time.Now().Add(-1*time.Minute).Unix())
	resetResp5, _ := h.ResetPassword(ctx, &user.ResetPasswordRequest{
		Email:       email,
		Token:       "expiredtoken",
		NewPassword: "validpassword",
	})
	if resetResp5.Success {
		t.Error("should not reset with expired token")
	}
}
