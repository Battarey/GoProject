package repository

import (
	"user-service/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(id string) (*model.User, error) {
	var user model.User
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	if err := r.db.Where("id = ?", uuidID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) DeleteUser(id string) error {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Delete(&model.User{}, "id = ?", uuidID).Error
}

func (r *UserRepository) ListUsers(offset, limit int) ([]model.User, error) {
	var users []model.User
	if err := r.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
