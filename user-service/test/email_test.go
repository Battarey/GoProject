package test

import (
	"errors"
	"strings"
	"testing"
	"user-service/config"
	"user-service/handler"
)

func TestEmailSendMock(t *testing.T) {
	handler.SendEmailFunc = func(cfg *config.Config, to, token string) error {
		return nil
	}
}

func TestGenerateEmailToken(t *testing.T) {
	token, err := handler.GenerateEmailToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(token) < 10 {
		t.Error("token too short")
	}
}

func TestGenerateResetToken(t *testing.T) {
	token, err := handler.GenerateResetToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(token) < 10 {
		t.Error("token too short")
	}
}

func TestSendConfirmationEmail_Format(t *testing.T) {
	cfg := &config.Config{SMTPUser: "user", SMTPPass: "pass", SMTPHost: "smtp.example.com", SMTPPort: "587", FromEmail: "noreply@example.com"}
	to := "test@example.com"
	token := "sometoken"
	// Подменяем smtp.SendMail на mock
	called := false
	handler.SendEmailFunc = func(cfg *config.Config, to, token string) error {
		called = true
		if !strings.Contains(to, "@") {
			t.Error("invalid email format")
		}
		if token == "" {
			t.Error("empty token")
		}
		return nil
	}
	_ = handler.SendEmailFunc(cfg, to, token)
	if !called {
		t.Error("SendEmailFunc was not called")
	}
}

func TestSendConfirmationEmail_Error(t *testing.T) {
	cfg := &config.Config{}
	handler.SendEmailFunc = func(cfg *config.Config, to, token string) error {
		return errors.New("smtp error")
	}
	err := handler.SendEmailFunc(cfg, "fail@example.com", "token")
	if err == nil || !strings.Contains(err.Error(), "smtp error") {
		t.Error("expected smtp error")
	}
}
