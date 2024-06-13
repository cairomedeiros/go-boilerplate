package service

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cairomedeiros/go-boilerplate/internal/entity"
	"github.com/cairomedeiros/go-boilerplate/internal/repository"
	"gorm.io/gorm"
)

type IUserService interface {
	RegisterUser(w http.ResponseWriter, r *http.Request) (entity.User, error)
	Login(w http.ResponseWriter, r *http.Request) (entity.User, error)
}

type UserService struct {
	userRepository repository.IUserRepository
}

func NewUserService(userRepo repository.IUserRepository) IUserService {
	return &UserService{userRepository: userRepo}
}

func (s *UserService) RegisterUser(w http.ResponseWriter, r *http.Request) (entity.User, error) {
	reqBody := entity.User{}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return entity.User{}, err
	}

	foundUser, err := s.userRepository.FindUserByEmail(reqBody)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.userRepository.Register(reqBody)
		}

		return entity.User{}, err
	}

	return foundUser, errors.New("user already exists in database")
}

func (s *UserService) Login(w http.ResponseWriter, r *http.Request) (entity.User, error) {
	reqBody := entity.User{}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return entity.User{}, err
	}

	foundUser, err := s.userRepository.FindUserByEmail(reqBody)
	if err != nil {
		return reqBody, err
	}

	return foundUser, nil
}
