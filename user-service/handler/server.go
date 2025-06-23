package handler

import (
	pb "user-service/proto"
	"user-service/repository"
	"user-service/security"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	Repo       *repository.UserRepository
	JwtService *security.JWTService
}
