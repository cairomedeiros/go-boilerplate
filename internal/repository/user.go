package repository

import (
	"github.com/cairomedeiros/go-boilerplate/config"
	"github.com/cairomedeiros/go-boilerplate/internal/entity"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Register(user entity.User) (entity.User, error)
	FindUserByEmail(user entity.User) (entity.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() IUserRepository {
	return &UserRepository{db: config.GetPostgreSQL()}
}

func (r *UserRepository) Register(user entity.User) (entity.User, error) {

	err := r.db.Create(&user).Error

	if err != nil {
		return entity.User{}, err
	}

	return user, nil

}

func (r *UserRepository) FindUserByEmail(user entity.User) (entity.User, error) {
	var existingUser entity.User
	err := r.db.Where("email = ?", user.Email).First(&existingUser)
	if err.Error != nil {
		return entity.User{}, err.Error
	}

	return existingUser, nil
}
