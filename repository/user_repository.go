package repository

import (
	"gorm.io/gorm"
	"hrm/domain"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (userRepository *UserRepository) Create(user *domain.User) error {
	return userRepository.db.Create(user).Error
}

func (userRepository *UserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := userRepository.db.Where("email = ?", email).First(&user).Error
	return &user, err
}
