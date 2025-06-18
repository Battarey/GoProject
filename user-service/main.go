package main

import (
	"log"
	"net"

	"user-service/config"
	"user-service/handler"
	"user-service/model"
	pb "user-service/proto"
	"user-service/repository"
	"user-service/security"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	// Автоматическая миграция
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	repo := repository.NewUserRepository(db)
	jwtService := security.NewJWTService(cfg.JWTSecret)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &handler.UserServer{
		Repo:       repo,
		JwtService: jwtService,
	})
	log.Println("user-service started on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
