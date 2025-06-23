package handler

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net/smtp"
	"regexp"
	"strings"
	"sync"
	"time"
	"user-service/config"
	"user-service/model"
	pb "user-service/proto"
	"user-service/repository"
	"user-service/security"

	"github.com/google/uuid"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	Repo       *repository.UserRepository
	JwtService *security.JWTService
}

var emailRegex = regexp.MustCompile(`^[\w._%+-]+@[\w.-]+\.[a-zA-Z]{2,}$`)
var allowedRoles = map[string]bool{"user": true, "admin": true}

func isValidEmail(email string) bool {
	// Простейшая проверка email (можно заменить на regexp)
	return len(email) >= 6 && len(email) <= 128 && strings.Contains(email, "@")
}

func isValidRole(role string) bool {
	return role == "user" || role == "admin"
}

func validateRegisterInput(req *pb.RegisterRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if req.Role != "" && !allowedRoles[req.Role] {
		return errors.New("invalid role")
	}
	return nil
}

func validateUpdateInput(req *pb.UpdateUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email")
	}
	if req.Role != "" && !allowedRoles[req.Role] {
		return errors.New("invalid role")
	}
	return nil
}

// Rate limiting (in-memory, на pet-проект)
type rateLimiter struct {
	mu      sync.Mutex
	buckets map[string][]int64 // ключ: email или IP, значения: unix timestamps
	limit   int
	window  int64 // в секундах
}

func newRateLimiter(limit int, windowSec int64) *rateLimiter {
	return &rateLimiter{
		buckets: make(map[string][]int64),
		limit:   limit,
		window:  windowSec,
	}
}

func (r *rateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().Unix()
	windowStart := now - r.window
	times := r.buckets[key]
	// Оставляем только те, что в окне
	var filtered []int64
	for _, t := range times {
		if t > windowStart {
			filtered = append(filtered, t)
		}
	}
	if len(filtered) >= r.limit {
		return false
	}
	filtered = append(filtered, now)
	r.buckets[key] = filtered
	return true
}

var regLimiter = newRateLimiter(5, 60)    // 5 регистраций в минуту на email
var loginLimiter = newRateLimiter(10, 60) // 10 логинов в минуту на email

// MOCK: Подмена отправки email для тестов
var SendEmailFunc = sendConfirmationEmail

func (s *UserServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if !regLimiter.Allow(strings.ToLower(req.Email)) {
		return nil, errors.New("too many registration attempts, try later")
	}
	if err := validateRegisterInput(req); err != nil {
		return nil, err
	}
	existing, err := s.Repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	hash := sha256.Sum256([]byte(req.Password))
	hashedPassword := hex.EncodeToString(hash[:])
	role := req.Role
	if role == "" {
		role = "user"
	}
	token, err := generateEmailToken()
	if err != nil {
		return nil, err
	}
	user := &model.User{
		ID:                     uuid.New(),
		Username:               req.Username,
		Email:                  req.Email,
		Password:               hashedPassword,
		Role:                   role,
		IsEmailConfirmed:       false,
		EmailConfirmationToken: token,
	}
	if err := s.Repo.CreateUser(user); err != nil {
		return nil, err
	}
	cfg := config.LoadConfig()
	_ = SendEmailFunc(cfg, user.Email, token) // Ошибку можно логировать
	return &pb.RegisterResponse{UserId: user.ID.String()}, nil
}

func (s *UserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if !loginLimiter.Allow(strings.ToLower(req.Email)) {
		return nil, errors.New("too many login attempts, try later")
	}
	user, err := s.Repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	hash := sha256.Sum256([]byte(req.Password))
	if user.Password != hex.EncodeToString(hash[:]) {
		return nil, errors.New("invalid password")
	}
	token, err := s.JwtService.GenerateToken(user.ID.String())
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{Token: token}, nil
}

func (s *UserServer) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	user, err := s.Repo.GetUserByID(req.UserId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return &pb.GetProfileResponse{
		UserId:   user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if err := validateUpdateInput(req); err != nil {
		return nil, err
	}
	user, err := s.Repo.GetUserByID(req.UserId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	user.Username = req.Username
	user.Email = req.Email
	user.Role = req.Role
	if err := s.Repo.UpdateUser(user); err != nil {
		return nil, err
	}
	return &pb.UpdateUserResponse{UserId: user.ID.String()}, nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := s.Repo.DeleteUser(req.UserId)
	if err != nil {
		return &pb.DeleteUserResponse{Success: false}, err
	}
	return &pb.DeleteUserResponse{Success: true}, nil
}

func (s *UserServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	offset := (int(req.Page) - 1) * int(req.PageSize)
	limit := int(req.PageSize)
	users, err := s.Repo.ListUsers(offset, limit)
	if err != nil {
		return nil, err
	}
	var userInfos []*pb.UserInfo
	for _, u := range users {
		userInfos = append(userInfos, &pb.UserInfo{
			UserId:   u.ID.String(),
			Username: u.Username,
			Email:    u.Email,
			Role:     u.Role,
		})
	}
	return &pb.ListUsersResponse{Users: userInfos, Total: int32(len(userInfos))}, nil
}

func (s *UserServer) ConfirmEmail(ctx context.Context, req *pb.ConfirmEmailRequest) (*pb.ConfirmEmailResponse, error) {
	user, err := s.Repo.GetUserByEmailAndToken(req.Email, req.Token)
	if err != nil {
		return &pb.ConfirmEmailResponse{Success: false, Message: "internal error"}, err
	}
	if user == nil {
		return &pb.ConfirmEmailResponse{Success: false, Message: "invalid token or email"}, nil
	}
	if user.IsEmailConfirmed {
		return &pb.ConfirmEmailResponse{Success: false, Message: "email already confirmed"}, nil
	}
	if err := s.Repo.ConfirmUserEmail(user); err != nil {
		return &pb.ConfirmEmailResponse{Success: false, Message: "failed to confirm email"}, err
	}
	return &pb.ConfirmEmailResponse{Success: true, Message: "email confirmed"}, nil
}

func (s *UserServer) RequestPasswordReset(ctx context.Context, req *pb.RequestPasswordResetRequest) (*pb.RequestPasswordResetResponse, error) {
	user, err := s.Repo.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		return &pb.RequestPasswordResetResponse{Success: false, Message: "user not found"}, nil
	}
	token, err := generateResetToken()
	if err != nil {
		return &pb.RequestPasswordResetResponse{Success: false, Message: "failed to generate token"}, err
	}
	expiresAt := time.Now().Add(30 * time.Minute).Unix()
	if err := s.Repo.SetPasswordResetToken(user.Email, token, expiresAt); err != nil {
		return &pb.RequestPasswordResetResponse{Success: false, Message: "failed to save token"}, err
	}
	cfg := config.LoadConfig()
	_ = sendPasswordResetEmail(cfg, user.Email, token)
	return &pb.RequestPasswordResetResponse{Success: true, Message: "reset email sent"}, nil
}

func (s *UserServer) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	user, err := s.Repo.GetUserByEmailAndResetToken(req.Email, req.Token)
	if err != nil || user == nil {
		return &pb.ResetPasswordResponse{Success: false, Message: "invalid token or email"}, nil
	}
	if user.PasswordResetExpiresAt < time.Now().Unix() {
		return &pb.ResetPasswordResponse{Success: false, Message: "token expired"}, nil
	}
	if len(req.NewPassword) < 6 {
		return &pb.ResetPasswordResponse{Success: false, Message: "password too short"}, nil
	}
	hash := sha256.Sum256([]byte(req.NewPassword))
	hashedPassword := hex.EncodeToString(hash[:])
	if err := s.Repo.ResetPassword(user, hashedPassword); err != nil {
		return &pb.ResetPasswordResponse{Success: false, Message: "failed to reset password"}, err
	}
	return &pb.ResetPasswordResponse{Success: true, Message: "password reset successful"}, nil
}

func generateEmailToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func sendConfirmationEmail(cfg *config.Config, to, token string) error {
	msg := fmt.Sprintf("Subject: Email Confirmation\n\nPlease confirm your email by clicking the link: http://localhost:8080/confirm?email=%s&token=%s", to, token)
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)
	return smtp.SendMail(addr, auth, cfg.FromEmail, []string{to}, []byte(msg))
}

func generateResetToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func sendPasswordResetEmail(cfg *config.Config, to, token string) error {
	msg := fmt.Sprintf("Subject: Password Reset\n\nTo reset your password, click the link: http://localhost:8080/reset?email=%s&token=%s", to, token)
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)
	return smtp.SendMail(addr, auth, cfg.FromEmail, []string{to}, []byte(msg))
}
