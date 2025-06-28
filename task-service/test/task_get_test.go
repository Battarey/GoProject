package test

import (
	"context"
	"task-service/proto"
	"testing"
)

func TestGetTask_NotFound(t *testing.T) {
	ts := setupTestServer(t)
	ctx := context.Background()
	_, err := ts.GetTask(ctx, &proto.GetTaskRequest{TaskId: "nonexistent-id"})
	if err == nil {
		t.Error("expected error for non-existent task, got nil")
	}
}
