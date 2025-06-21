package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
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

func (s *UserServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	existing, err := s.Repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	hash := sha256.Sum256([]byte(req.Password))
	hashedPassword := hex.EncodeToString(hash[:])
	user := &model.User{
		ID:       uuid.New(),
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}
	if err := s.Repo.CreateUser(user); err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{UserId: user.ID.String()}, nil
}

func (s *UserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
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
	}, nil
}
