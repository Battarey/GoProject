package main

import (
	"context"
	"os"
	"testing"
	"time"

	taskpb "task-service/proto"
	userpb "user-service/proto"

	"google.golang.org/grpc"
)

func TestE2E_UserAndTaskFlow(t *testing.T) {
	userConn, err := grpc.Dial("user-service:50051", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		t.Fatalf("failed to connect to user-service: %v", err)
	}
	defer userConn.Close()
	userClient := userpb.NewUserServiceClient(userConn)

	taskConn, err := grpc.Dial("task-service:50052", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		t.Fatalf("failed to connect to task-service: %v", err)
	}
	defer taskConn.Close()
	taskClient := taskpb.NewTaskServiceClient(taskConn)

	ctx := context.Background()

	// 1. Создать пользователя
	userResp, err := userClient.Register(ctx, &userpb.RegisterRequest{
		Username: "e2euser",
		Password: "e2epass",
		Email:    "e2e@example.com",
	})
	if err != nil {
		t.Fatalf("user register failed: %v", err)
	}

	// 2. Войти и получить JWT
	tokenResp, err := userClient.Login(ctx, &userpb.LoginRequest{
		Username: "e2euser",
		Password: "e2epass",
	})
	if err != nil {
		t.Fatalf("user login failed: %v", err)
	}
	jwt := tokenResp.Token

	// 3. Создать задачу с этим JWT
	md := map[string][]string{"authorization": {"Bearer " + jwt}}
	ctxWithJWT := grpc.NewOutgoingContext(context.Background(), md)
	taskResp, err := taskClient.CreateTask(ctxWithJWT, &taskpb.CreateTaskRequest{Title: "E2E task"})
	if err != nil {
		t.Fatalf("create task failed: %v", err)
	}
	if taskResp.TaskId == "" {
		t.Fatalf("empty task id")
	}
}

func main() {
	os.Exit(m.Run())
}
