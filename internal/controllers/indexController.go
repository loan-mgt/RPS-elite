package controllers

import (
	"net/http"

	"rcp/elite/internal/utils"
)

func IndexController(w http.ResponseWriter, r *http.Request) {
	err := utils.Templates.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
