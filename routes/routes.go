package router

import (
	"github.com/cairomedeiros/go-boilerplate/internal/controller"
	"github.com/cairomedeiros/go-boilerplate/internal/repository"
	"github.com/cairomedeiros/go-boilerplate/internal/service"
	"github.com/gorilla/mux"
)

func initializeRoutes(r *mux.Router) {

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository)
	userControler := controller.NewUserController(userService)

	r.HandleFunc("/register", userControler.RegisterUser).Methods("POST")
	r.HandleFunc("/login", userControler.Login).Methods("POST")
}
