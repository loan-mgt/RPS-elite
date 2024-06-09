package routes

import (
	"net/http"
	"rcp/elite/internal/controllers"

	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {

	r.HandleFunc("/ws", controllers.MainController)
	r.HandleFunc("/", controllers.IndexController)

	r.PathPrefix("/assets/").Handler(http.FileServer(http.Dir("static/")))
	r.PathPrefix("/styles/").Handler(http.FileServer(http.Dir("static/")))
	r.PathPrefix("/scripts/").Handler(http.FileServer(http.Dir("static/")))
}
