package routes

import (
	"net/http"
	"rcp/elite/internal/controllers"

	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {

	r.HandleFunc("/ws", controllers.MainController)
	r.HandleFunc("/", controllers.IndexController)

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("static/assets"))))
	r.PathPrefix("/styles/").Handler(http.StripPrefix("/styles/", http.FileServer(http.Dir("static/styles"))))
}
