package test

import (
	"user-service/handler"
	"user-service/model"
	"user-service/repository"
	"user-service/security"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupHandlerTest() *handler.UserServer {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.User{})
	repo := repository.NewUserRepository(db)
	jwt := security.NewJWTService("testsecret")
	return &handler.UserServer{Repo: repo, JwtService: jwt}
}
