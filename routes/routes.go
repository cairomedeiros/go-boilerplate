package router

import (
	"github.com/cairomedeiros/go-boilerplate/internal/controller"
	"github.com/cairomedeiros/go-boilerplate/internal/repository"
	"github.com/cairomedeiros/go-boilerplate/internal/service"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

func initializeRoutes(r *mux.Router) {

	//Initialize validator
	validate := validator.New()
	auth := service.NewAuthHandlerImpl(validate)

	// Protected Routes
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(auth.AuthMiddleware)

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, auth)
	userControler := controller.NewUserController(userService)

	r.HandleFunc("/register", userControler.RegisterUser).Methods("POST")
	r.HandleFunc("/login", userControler.Login).Methods("POST")
}
