package controllers

import (
	"net/http"
	"text/template"
)

var templates = template.Must(template.ParseGlob("internal/component/*.gohtml"))

func IndexController(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
