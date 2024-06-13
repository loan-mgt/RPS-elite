package controllers

import (
	"log"
	"net/http"

	"rcp/elite/internal/utils"
)

func IndexController(w http.ResponseWriter, r *http.Request) {
	// Check if the 'stk' cookie is set
	cookie, err := r.Cookie("stk")
	if err != nil {
		err = utils.Templates.ExecuteTemplate(w, "index", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	// If the cookie is found, you can access its value with cookie.Value
	log.Println("Cookie 'stk' found with value:", cookie.Value)

	err = utils.Templates.ExecuteTemplate(w, "logged", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
}
