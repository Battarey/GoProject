package test

import (
	"testing"
	"user-service/model"
	"user-service/repository"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateAndGetUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	db.AutoMigrate(&model.User{})

	repo := repository.NewUserRepository(db)
	user := &model.User{
		ID:       "test-id",
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpass",
	}
	err = repo.CreateUser(user)
	require.NoError(t, err)

	got, err := repo.GetUserByEmail("test@example.com")
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, user.ID, got.ID)
}
