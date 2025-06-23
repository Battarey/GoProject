package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"user-service/config"
)

var SendEmailFunc = SendConfirmationEmail

func SendConfirmationEmail(cfg *config.Config, to, token string) error {
	msg := fmt.Sprintf("Subject: Email Confirmation\n\nPlease confirm your email by clicking the link: http://localhost:8080/confirm?email=%s&token=%s", to, token)
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)
	return smtp.SendMail(addr, auth, cfg.FromEmail, []string{to}, []byte(msg))
}

func SendPasswordResetEmail(cfg *config.Config, to, token string) error {
	msg := fmt.Sprintf("Subject: Password Reset\n\nTo reset your password, click the link: http://localhost:8080/reset?email=%s&token=%s", to, token)
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)
	return smtp.SendMail(addr, auth, cfg.FromEmail, []string{to}, []byte(msg))
}

func GenerateEmailToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func GenerateResetToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
