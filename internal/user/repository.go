package user

import (
	"fmt"
	"go/adv-demo/pkg/db"
)

type UserRepository struct {
	db *db.Db
}

func NewUserRepository(db *db.Db) *UserRepository {
	return &UserRepository{db}
}

func (repo UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	result := repo.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, result.Error

}

func (repo UserRepository) Create(user *User) (*User, error) {
	result := repo.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
func (repo *UserRepository) FindUserByName(username string) (*User, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (repo *UserRepository) FindUserById(id int) (*User, error) {
	return nil, fmt.Errorf("Not implemented")
}
