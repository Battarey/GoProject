package test

import (
	"testing"
	"user-service/config"
	"user-service/handler"
)

func TestEmailSendMock(t *testing.T) {
	handler.SendEmailFunc = func(cfg *config.Config, to, token string) error {
		return nil
	}
	// Можно добавить дополнительные проверки/mock-тесты для email-отправки
}
