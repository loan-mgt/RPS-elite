package routes

import (
	"github.com/gorilla/mux"

	"rcp/elite/internal/controllers"
)

// SetupRoutes configures the routes for the application
func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/ws", controllers.MainController)
}
