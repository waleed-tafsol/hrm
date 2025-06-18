package usecase

import (
	"golang.org/x/crypto/bcrypt"
	"hrm/domain"
)

type UserService struct {
	userRepository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (userService *UserService) SignUp(user *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return userService.userRepository.Create(user)
}
