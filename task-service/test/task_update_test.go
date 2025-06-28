package test

import (
	"context"
	"task-service/proto"
	"testing"
)

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
