package test

import (
	"context"
	"task-service/handler"
	"task-service/model"
	"task-service/proto"
	"task-service/repository"
	"task-service/security"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
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
	rateLimiter := handler.NewRateLimiter(10 * time.Millisecond) // минимальный лимит для тестов
	return &handler.TaskServer{Repo: repo, JwtService: jwtService, RateLimiter: rateLimiter}
}

func makeJWT(t *testing.T, secret, userID, role string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": userID, "role": role})
	tokStr, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign jwt: %v", err)
	}
	return tokStr
}

func ctxWithJWT(token string) context.Context {
	md := metadata.New(map[string]string{"authorization": "Bearer " + token})
	return metadata.NewIncomingContext(context.Background(), md)
}

func TestCreateTask(t *testing.T) {
	ts := setupTestServer(t)
	secret := "testsecret"
	userUUID := "11111111-1111-1111-1111-111111111111"
	jwt1 := makeJWT(t, secret, userUUID, "user")
	ctx := ctxWithJWT(jwt1)
	resp, err := ts.CreateTask(ctx, &proto.CreateTaskRequest{Title: "Test task"})
	if err != nil || resp.TaskId == "" {
		t.Errorf("expected task to be created, got err: %v", err)
	}
}

func TestGetTask_NotFound(t *testing.T) {
	ts := setupTestServer(t)
	ctx := context.Background()
	_, err := ts.GetTask(ctx, &proto.GetTaskRequest{TaskId: "nonexistent-id"})
	if err == nil {
		t.Error("expected error for non-existent task, got nil")
	}
}

func TestUpdateTask(t *testing.T) {
	ts := setupTestServer(t)
	secret := "testsecret"
	userUUID := "11111111-1111-1111-1111-111111111111"
	jwt1 := makeJWT(t, secret, userUUID, "user")
	ctx := ctxWithJWT(jwt1)
	createResp, _ := ts.CreateTask(ctx, &proto.CreateTaskRequest{Title: "To update"})
	_, err := ts.UpdateTask(ctx, &proto.UpdateTaskRequest{TaskId: createResp.TaskId, Title: "Updated"})
	if err != nil {
		t.Errorf("expected update to succeed, got err: %v", err)
	}
}

func TestDeleteTask(t *testing.T) {
	ts := setupTestServer(t)
	secret := "testsecret"
	userUUID := "11111111-1111-1111-1111-111111111111"
	jwt1 := makeJWT(t, secret, userUUID, "user")
	ctx := ctxWithJWT(jwt1)
	createResp, _ := ts.CreateTask(ctx, &proto.CreateTaskRequest{Title: "To delete"})
	_, err := ts.DeleteTask(ctx, &proto.DeleteTaskRequest{TaskId: createResp.TaskId})
	if err != nil {
		t.Errorf("expected delete to succeed, got err: %v", err)
	}
}

func TestCreateTask_Validation(t *testing.T) {
	ts := setupTestServer(t)
	secret := "testsecret"
	userUUID := "11111111-1111-1111-1111-111111111111"
	jwt1 := makeJWT(t, secret, userUUID, "user")
	ctx := ctxWithJWT(jwt1)
	_, err := ts.CreateTask(ctx, &proto.CreateTaskRequest{Title: ""})
	if err == nil {
		t.Error("expected validation error for empty title")
	}
}

func TestChangeStatus(t *testing.T) {
	ts := setupTestServer(t)
	secret := "testsecret"
	userUUID := "11111111-1111-1111-1111-111111111111"
	jwt1 := makeJWT(t, secret, userUUID, "user")
	ctx := ctxWithJWT(jwt1)
	createResp, _ := ts.CreateTask(ctx, &proto.CreateTaskRequest{Title: "Status test"})
	_, err := ts.ChangeStatus(ctx, &proto.ChangeStatusRequest{TaskId: createResp.TaskId, Status: "done"})
	if err != nil {
		t.Errorf("expected status change to succeed, got err: %v", err)
	}
}

func TestChangeStatus_NotFound(t *testing.T) {
	ts := setupTestServer(t)
	secret := "testsecret"
	userUUID := "11111111-1111-1111-1111-111111111111"
	jwt1 := makeJWT(t, secret, userUUID, "user")
	ctx := ctxWithJWT(jwt1)
	_, err := ts.ChangeStatus(ctx, &proto.ChangeStatusRequest{TaskId: "nonexistent-id", Status: "done"})
	if err == nil {
		t.Error("expected error for non-existent task on status change")
	}
}

func TestUpdateTask_Auth_EdgeCases(t *testing.T) {
	ts := setupTestServer(t)
	secret := "testsecret"
	user1 := "11111111-1111-1111-1111-111111111111"
	jwt1 := makeJWT(t, secret, user1, "user")
	ctx1 := ctxWithJWT(jwt1)
	resp, err := ts.CreateTask(ctx1, &proto.CreateTaskRequest{Title: "Edge task"})
	if err != nil {
		t.Fatalf("create task failed: %v", err)
	}
	taskID := resp.TaskId

	// Попытка обновить без токена
	_, err = ts.UpdateTask(context.Background(), &proto.UpdateTaskRequest{TaskId: taskID, Title: "No token"})
	if err == nil || err.Error() == "" {
		t.Error("expected error for missing token")
	}

	// Попытка с невалидным токеном
	ctxBad := ctxWithJWT("bad.token.value")
	_, err = ts.UpdateTask(ctxBad, &proto.UpdateTaskRequest{TaskId: taskID, Title: "Bad token"})
	if err == nil || err.Error() == "" {
		t.Error("expected error for invalid token")
	}

	// Попытка чужим user_id
	jwt2 := makeJWT(t, secret, "22222222-2222-2222-2222-222222222222", "user")
	ctx2 := ctxWithJWT(jwt2)
	_, err = ts.UpdateTask(ctx2, &proto.UpdateTaskRequest{TaskId: taskID, Title: "Not owner"})
	if err == nil || err.Error() == "" || err.Error() == "unauthorized" {
		t.Error("expected permission denied for not owner")
	}

	// Успешно своим user_id
	_, err = ts.UpdateTask(ctx1, &proto.UpdateTaskRequest{TaskId: taskID, Title: "Owner update"})
	if err != nil {
		t.Errorf("expected update by owner, got err: %v", err)
	}

	// Успешно с ролью admin
	jwtAdmin := makeJWT(t, secret, "admin-id", "admin")
	ctxAdmin := ctxWithJWT(jwtAdmin)
	_, err = ts.UpdateTask(ctxAdmin, &proto.UpdateTaskRequest{TaskId: taskID, Title: "Admin update"})
	if err != nil {
		t.Errorf("expected update by admin, got err: %v", err)
	}
}

func TestDeleteTask_Auth_EdgeCases(t *testing.T) {
	ts := setupTestServer(t)
	secret := "testsecret"
	user1 := "11111111-1111-1111-1111-111111111111"
	jwt1 := makeJWT(t, secret, user1, "user")
	ctx1 := ctxWithJWT(jwt1)
	resp, err := ts.CreateTask(ctx1, &proto.CreateTaskRequest{Title: "Edge del"})
	if err != nil {
		t.Fatalf("create task failed: %v", err)
	}
	taskID := resp.TaskId

	// Попытка удалить без токена
	_, err = ts.DeleteTask(context.Background(), &proto.DeleteTaskRequest{TaskId: taskID})
	if err == nil || err.Error() == "" {
		t.Error("expected error for missing token on delete")
	}

	// Попытка с невалидным токеном
	ctxBad := ctxWithJWT("bad.token.value")
	_, err = ts.DeleteTask(ctxBad, &proto.DeleteTaskRequest{TaskId: taskID})
	if err == nil || err.Error() == "" {
		t.Error("expected error for invalid token on delete")
	}

	// Попытка чужим user_id
	jwt2 := makeJWT(t, secret, "22222222-2222-2222-2222-222222222222", "user")
	ctx2 := ctxWithJWT(jwt2)
	_, err = ts.DeleteTask(ctx2, &proto.DeleteTaskRequest{TaskId: taskID})
	if err == nil || err.Error() == "" || err.Error() == "unauthorized" {
		t.Error("expected permission denied for not creator")
	}

	// Успешно своим user_id (creator)
	_, err = ts.DeleteTask(ctx1, &proto.DeleteTaskRequest{TaskId: taskID})
	if err != nil {
		t.Errorf("expected delete by creator, got err: %v", err)
	}

	// Создаём новую задачу для admin
	resp, err = ts.CreateTask(ctx1, &proto.CreateTaskRequest{Title: "Admin del"})
	taskID = resp.TaskId
	jwtAdmin := makeJWT(t, secret, "admin-id", "admin")
	ctxAdmin := ctxWithJWT(jwtAdmin)
	_, err = ts.DeleteTask(ctxAdmin, &proto.DeleteTaskRequest{TaskId: taskID})
	if err != nil {
		t.Errorf("expected delete by admin, got err: %v", err)
	}
}
