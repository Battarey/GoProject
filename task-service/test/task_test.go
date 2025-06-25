package test

import (
	"context"
	"task-service/handler"
	"task-service/model"
	"task-service/proto"
	"task-service/repository"
	"task-service/security"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestServer(t *testing.T) *handler.TaskServer {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	db.AutoMigrate(&model.Task{})
	repo := repository.NewTaskRepository(db)
	jwtService := security.NewJWTService("testsecret")
	return &handler.TaskServer{Repo: repo, JwtService: jwtService}
}

func TestCreateTask(t *testing.T) {
	ts := setupTestServer(t)
	ctx := context.Background()
	resp, err := ts.CreateTask(ctx, &proto.CreateTaskRequest{Title: "Test task"})
	if err != nil || resp.TaskId == "" {
		t.Errorf("expected task to be created, got err: %v", err)
	}
}
