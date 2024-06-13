package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cairomedeiros/go-boilerplate/config"
	"github.com/cairomedeiros/go-boilerplate/helper"
	"github.com/cairomedeiros/go-boilerplate/internal/entity"
	"github.com/cairomedeiros/go-boilerplate/internal/repository"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type IUserService interface {
	RegisterUser(w http.ResponseWriter, r *http.Request) (entity.User, error)
	Login(w http.ResponseWriter, r *http.Request) (entity.User, error)
}

type UserService struct {
	userRepository repository.IUserRepository
	authHandler    *AuthHandler
}

type AuthHandler struct {
	Validate *validator.Validate
}

func NewUserService(userRepo repository.IUserRepository, authHandler *AuthHandler) IUserService {
	return &UserService{
		userRepository: userRepo,
		authHandler:    authHandler,
	}
}

func NewAuthHandlerImpl(validate *validator.Validate) *AuthHandler {
	return &AuthHandler{Validate: validate}
}

func (s *UserService) RegisterUser(w http.ResponseWriter, r *http.Request) (entity.User, error) {
	reqBody := entity.User{}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return entity.User{}, err
	}

	if err := s.authHandler.validateUser(reqBody); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessage := fmt.Sprintf("Validation failed for field: %s", validationErrors[0].Field())
		return reqBody, errors.New(errorMessage)
	}

	foundUser, err := s.userRepository.FindUserByEmail(reqBody)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			password, err := helper.EncryptPassword(reqBody.Password)

			if err != nil {

				return reqBody, err
			}

			newUser := entity.User{
				Name:     reqBody.Name,
				Email:    reqBody.Email,
				Password: password,
			}

			return s.userRepository.Register(newUser)
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

	if err := s.authHandler.validateUser(reqBody); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessage := fmt.Sprintf("Validation failed for field: %s", validationErrors[0].Field())
		return reqBody, errors.New(errorMessage)
	}

	foundUser, err := s.userRepository.FindUserByEmail(reqBody)
	if err != nil {
		return reqBody, err
	}

	valid := helper.ComparePassword(reqBody.Password, foundUser.Password)

	if !valid {
		return reqBody, errors.New("password invalid")
	}

	token, err := helper.CreateToken(foundUser.Email)

	if err != nil {
		return reqBody, errors.New("jwt error")
	}

	fmt.Print(token) //adjust

	return foundUser, nil
}

func (a *AuthHandler) validateUser(reqBody entity.User) error {
	if err := a.Validate.Struct(reqBody); err != nil {
		return err
	}
	return nil
}

func (c AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return config.SecretKey, nil
		})

		if err != nil {
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "email", claims["email"])
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
