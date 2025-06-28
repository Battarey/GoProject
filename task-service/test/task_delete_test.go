package test

import (
	"context"
	"task-service/proto"
	"testing"
)

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
