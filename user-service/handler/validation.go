package handler

import (
	"errors"
	"regexp"
	"strings"
	pb "user-service/proto"
)

var EmailRegex = regexp.MustCompile(`^[\w._%+-]+@[\w.-]+\.[a-zA-Z]{2,}$`)
var AllowedRoles = map[string]bool{"user": true, "admin": true}

func IsValidEmail(email string) bool {
	// Простейшая проверка email (можно заменить на regexp)
	return len(email) >= 6 && len(email) <= 128 && strings.Contains(email, "@")
}

func IsValidRole(role string) bool {
	return role == "user" || role == "admin"
}

func ValidateRegisterInput(req *pb.RegisterRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if !EmailRegex.MatchString(req.Email) {
		return errors.New("invalid email")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if req.Role != "" && !AllowedRoles[req.Role] {
		return errors.New("invalid role")
	}
	return nil
}

func ValidateUpdateInput(req *pb.UpdateUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if !EmailRegex.MatchString(req.Email) {
		return errors.New("invalid email")
	}
	if req.Role != "" && !AllowedRoles[req.Role] {
		return errors.New("invalid role")
	}
	return nil
}
