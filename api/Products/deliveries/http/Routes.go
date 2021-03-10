package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrdhira/warpin-test/api/Products/deliveries/http/controllers"
)

// Route struct
type Route struct{}

// Init func
func (r *Route) Init() *mux.Router {
	// Initialize Controllers
	productsControllers := controllers.InitProductsControllers()

	// Initialize Router
	Router := mux.NewRouter().StrictSlash(true)

	// Products Routes with no Auth
	ProductsNoAuthRoutes := Router.PathPrefix("/products").Subrouter()
	ProductsNoAuthRoutes.HandleFunc("/", productsControllers.GetProducts).Methods(http.MethodGet)
	ProductsNoAuthRoutes.HandleFunc("/{id}", productsControllers.GetProductsByID).Methods(http.MethodGet)

	// Products Routes with Auth Admin
	ProductsAuthAdminRoutes := Router.PathPrefix("/products/internal").Subrouter()
	ProductsAuthAdminRoutes.Use(AuthAdmniMiddleware)
	ProductsAuthAdminRoutes.HandleFunc("/add", productsControllers.AddProducts).Methods(http.MethodPost)
	ProductsAuthAdminRoutes.HandleFunc("/{id}", productsControllers.UpdateProducts).Methods(http.MethodPut)

	return Router
}
