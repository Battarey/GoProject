package main

import (
	"log"
	"net"

	"task-service/config"
	"task-service/handler"
	"task-service/proto"
	"task-service/repository"
	"task-service/security"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()

	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	repo := repository.NewTaskRepository(db)
	jwtService := security.NewJWTService(cfg.JWTSecret)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterTaskServiceServer(s, &handler.TaskServer{
		Repo:       repo,
		JwtService: jwtService,
	})

	log.Printf("task-service started on :%s", cfg.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
