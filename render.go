package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cbroglie/mustache"
)

func Render(w http.ResponseWriter, file_name_look_up string, template_data ...interface{}) {
	if file_name_look_up == "" {
		file_name_look_up = "views/pages/home.mst"
	}
	if file_name_look_up == "404page" {
		file_name_look_up = "views/pages/not-found.mst"
	}

	template, err := viewsDir.ReadFile(file_name_look_up)
	log.Println(file_name_look_up)
	if err != nil {
		log.Println("Error reading file", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	parsedTemplate, _ := mustache.ParseString(string(template))
	content, err := parsedTemplate.Render(template_data)
	if err != nil {
		log.Println("Error rendering template", err)
		http.Error(w, "Internal server error two", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, content)
}
