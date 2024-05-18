package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"rcp/elite/internal/routes"
)

func main() {
	r := mux.NewRouter()

	routes.SetupRoutes(r)

	http.Handle("/", r)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
