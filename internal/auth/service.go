package auth

import (
	"errors"
	"go/adv-demo/internal/user"
	"go/adv-demo/pkg/di"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo di.IUserRepository
}

func NewAuthService(userRepo di.IUserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (service *AuthService) Register(email, username, password string) (string, error) {
	existingUser, _ := service.userRepo.FindByEmail(email)
	if existingUser != nil {
		return "", errors.New(ErrUserExists)
	}
	hasheDPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &user.User{
		Email:    email,
		Password: string(hasheDPass),
		Name:     username,
	}
	_, err = service.userRepo.Create(user)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}

func (service *AuthService) Login(email, password string) (string, error) {
	user, err := service.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New(ErrWrongCredentials)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New(ErrWrongCredentials)
	}
	return user.Email, nil
}
