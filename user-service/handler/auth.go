package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"
	"user-service/config"
	"user-service/model"
	pb "user-service/proto"

	"github.com/google/uuid"
)

func (s *UserServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if !RegLimiter.Allow(strings.ToLower(req.Email)) {
		return nil, errors.New("too many registration attempts, try later")
	}
	if err := ValidateRegisterInput(req); err != nil {
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
	token, err := GenerateEmailToken()
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
	if !LoginLimiter.Allow(strings.ToLower(req.Email)) {
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
	token, err := GenerateResetToken()
	if err != nil {
		return &pb.RequestPasswordResetResponse{Success: false, Message: "failed to generate token"}, err
	}
	expiresAt := time.Now().Add(30 * time.Minute).Unix()
	if err := s.Repo.SetPasswordResetToken(user.Email, token, expiresAt); err != nil {
		return &pb.RequestPasswordResetResponse{Success: false, Message: "failed to save token"}, err
	}
	cfg := config.LoadConfig()
	_ = SendPasswordResetEmail(cfg, user.Email, token)
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
