package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrdhira/warpin-test/api/Orders/deliveries/http/controllers"
)

// Route struct
type Route struct{}

// Init func
func (r *Route) Init() *mux.Router {
	// Initialize Controllers
	ordersControllers := controllers.InitOrdersControllers()

	// Initialize Router
	Router := mux.NewRouter().StrictSlash(true)

	// Orders Routes with Auth
	OrdersAuthRoutes := Router.PathPrefix("/orders").Subrouter()
	OrdersAuthRoutes.Use(AuthMiddleware)
	OrdersAuthRoutes.HandleFunc("/", ordersControllers.OrdersListUsers).Methods(http.MethodGet)
	OrdersAuthRoutes.HandleFunc("/", ordersControllers.OrdersCreate).Methods(http.MethodPost)
	OrdersAuthRoutes.HandleFunc("/{id}", ordersControllers.OrdersUpdate).Methods(http.MethodPut)
	OrdersAuthRoutes.HandleFunc("/{id}/cancel", ordersControllers.OrdersCancel).Methods(http.MethodPut)

	// Users Routes with Auth Admin
	OrdersAuthAdminRoutes := Router.PathPrefix("/orders/internal").Subrouter()
	OrdersAuthAdminRoutes.Use(AuthAdmniMiddleware)
	OrdersAuthAdminRoutes.HandleFunc("/", ordersControllers.OrdersListAdmin).Methods(http.MethodGet)
	OrdersAuthAdminRoutes.HandleFunc("/{id}/approve", ordersControllers.OrdersApprove).Methods(http.MethodPut)
	OrdersAuthAdminRoutes.HandleFunc("/{id}/reject", ordersControllers.OrdersReject).Methods(http.MethodPut)

	return Router
}
