package controllers

import (
	"log"
	"net/http"

	templatedata "rcp/elite/internal/types/template-data"
	"rcp/elite/internal/utils"
)

func IndexController(w http.ResponseWriter, r *http.Request) {

	data := templatedata.IndexData{
		Main: "logged",
	}
	cookie, err := r.Cookie("stk")
	if err != nil {
		data.Main = "home"

		defaultCookie := http.Cookie{
			Name:     "stk",
			Value:    "default_value_here",
			MaxAge:   60 * 60 * 24 * 60, // 2 months in seconds (60 seconds * 60 minutes * 24 hours * 60 days)
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, &defaultCookie)
		log.Println("Cookie 'stk' set with default value")

	} else {
		// If the cookie is found, you can access its value with cookie.Value
		log.Println("Cookie 'stk' found with value:", cookie.Value)
	}

	err = utils.Templates.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
