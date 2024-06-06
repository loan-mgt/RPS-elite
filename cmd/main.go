package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"rcp/elite/internal/routes"
	"rcp/elite/internal/services"
	"rcp/elite/internal/utils"
)

func main() {
	utils.InitializeTemplates("template/*.tmpl")

	r := mux.NewRouter()

	routes.SetupRoutes(r)

	http.Handle("/", r)

	go services.SearchForGameToCreate()
	go services.StartPollMonitor()

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
