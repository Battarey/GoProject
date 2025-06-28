package test

import (
	"task-service/proto"
	"testing"
)

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
