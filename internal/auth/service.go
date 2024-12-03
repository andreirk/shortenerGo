package auth

import (
	"errors"
	"go/adv-demo/internal/user"
)

type AuthService struct {
	userRepo *user.UserRepository
}

func NewAuthService(userRepo *user.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (service *AuthService) Register(email, username, password string) (string, error) {
	existingUser, _ := service.userRepo.FindByEmail(email)
	if existingUser != nil {
		return "", errors.New(ErrUserExists)
	}
	user := &user.User{
		Email:    email,
		Password: "password",
		Name:     username,
	}
	_, err := service.userRepo.Create(user)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}
