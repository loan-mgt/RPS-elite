package utils

import (
	"html/template"
	"log"
)

var Templates *template.Template

func InitializeTemplates(pattern string) {
	var err error
	Templates, err = template.ParseGlob(pattern)
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
}
