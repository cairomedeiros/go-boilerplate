package repository

import (
	"fmt"

	"github.com/cairomedeiros/go-boilerplate/config"
	"github.com/cairomedeiros/go-boilerplate/internal/entity"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Register(user entity.User) (entity.User, error)
	FindUserByID(user entity.User) (entity.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() IUserRepository {
	return &UserRepository{db: config.GetPostgreSQL()}
}

func (r *UserRepository) Register(user entity.User) (entity.User, error) {

	newUser := entity.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	err := r.db.Create(&newUser).Error

	if err != nil {
		return entity.User{}, err
	}

	return newUser, nil

}

func (r *UserRepository) FindUserByID(user entity.User) (entity.User, error) {
	var existingUser entity.User
	err := r.db.Where("email = ?", user.Email).First(&existingUser)
	if err.Error != nil {
		fmt.Println("entrou no if do finduser: err", err)

		return entity.User{}, err.Error
	}
	fmt.Println("debugar isso direito caireca", existingUser)
	fmt.Println("finduser: err", err.Error)

	return existingUser, nil
}
