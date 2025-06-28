package test

import (
	"task-service/proto"
	"testing"
)

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
