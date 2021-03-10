package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrdhira/warpin-test/api/Users/deliveries/http/controllers"
)

// Route struct
type Route struct{}

// Init func
func (r *Route) Init() *mux.Router {
	// Initialize Controllers
	usersControllers := controllers.InitUsersControllers()

	// Initialize Router
	Router := mux.NewRouter().StrictSlash(true)

	// Users Routes with no Auth
	UsersNoAuthRoutes := Router.PathPrefix("/users").Subrouter()
	UsersNoAuthRoutes.HandleFunc("/register", usersControllers.Register).Methods(http.MethodPost)
	UsersNoAuthRoutes.HandleFunc("/login", usersControllers.Login).Methods(http.MethodPost)

	// Users Routes with Auth
	UsersAuthRoutes := Router.PathPrefix("/users").Subrouter()
	UsersAuthRoutes.Use(AuthMiddleware)
	UsersAuthRoutes.HandleFunc("/profile", usersControllers.Profile).Methods(http.MethodGet)
	UsersAuthRoutes.HandleFunc("/update-profile", usersControllers.UpdateProfile).Methods(http.MethodPut)
	UsersAuthRoutes.HandleFunc("/update-password", usersControllers.UpdatePassword).Methods(http.MethodPut)

	return Router
}
