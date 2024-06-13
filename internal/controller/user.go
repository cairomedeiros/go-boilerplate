package controller

import (
	"encoding/json"
	"net/http"

	"github.com/cairomedeiros/go-boilerplate/internal/service"
)

type UserController struct {
	UserService service.IUserService
}

func NewUserController(user service.IUserService) *UserController {
	return &UserController{UserService: user}
}

func (s *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	user, err := s.UserService.RegisterUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
	json.NewEncoder(w).Encode(user)
}

func (s *UserController) Login(w http.ResponseWriter, r *http.Request) {
	user, err := s.UserService.Login(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User logged!"))
	json.NewEncoder(w).Encode(user)
}
