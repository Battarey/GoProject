package handler

import (
	pb "task-service/proto"
	"task-service/repository"
	"task-service/security"
)

type TaskServer struct {
	pb.UnimplementedTaskServiceServer
	Repo       *repository.TaskRepository
	JwtService *security.JWTService
}
