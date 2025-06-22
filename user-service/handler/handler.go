package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
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

func (s *UserServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
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
	user := &model.User{
		ID:       uuid.New(),
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     role,
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
