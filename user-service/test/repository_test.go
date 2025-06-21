package test

import (
	"testing"
	"user-service/model"
	"user-service/repository"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect test db")
	}
	db.AutoMigrate(&model.User{})
	return db
}

func TestCreateAndGetUser(t *testing.T) {
	db := setupTestDB()
	repo := repository.NewUserRepository(db)

	user := &model.User{
		ID:       uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	err := repo.CreateUser(user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	got, err := repo.GetUserByEmail("test@example.com")
	if err != nil {
		t.Fatalf("failed to get user by email: %v", err)
	}
	if got.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, got.Email)
	}
}
