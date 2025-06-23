package test

import (
	"context"
	"strconv"
	"testing"
	user "user-service/proto"
)

func TestUpdateAndDeleteUser(t *testing.T) {
	h := SetupHandlerTest()
	ctx := context.Background()

	// Create a user
	regResp, err := h.Register(ctx, &user.RegisterRequest{
		Username: "John Doe",
		Email:    "john.doe@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Update the user's name
	updReq := &user.UpdateUserRequest{
		UserId:   regResp.UserId,
		Username: "Jane Doe",
		Email:    "john.doe@example.com",
		Role:     "user",
	}
	_, err = h.UpdateUser(ctx, updReq)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}

	// Delete the user
	delReq := &user.DeleteUserRequest{UserId: regResp.UserId}
	_, err = h.DeleteUser(ctx, delReq)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}
}

func TestListUsers(t *testing.T) {
	h := SetupHandlerTest()
	ctx := context.Background()

	// Create test users
	for i := 1; i <= 3; i++ {
		_, err := h.Register(ctx, &user.RegisterRequest{
			Username: "User" + strconv.Itoa(i),
			Email:    "user" + strconv.Itoa(i) + "@example.com",
			Password: "password123",
			Role:     "user",
		})
		if err != nil {
			t.Fatalf("register failed: %v", err)
		}
	}

	// List users
	req := &user.ListUsersRequest{Page: 1, PageSize: 10}
	res, err := h.ListUsers(ctx, req)
	if err != nil {
		t.Fatalf("list users failed: %v", err)
	}

	if len(res.Users) != 3 {
		t.Fatalf("expected 3 users, got %d", len(res.Users))
	}
}
