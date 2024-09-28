package main

import (
	"log"
	"net/http"
)

func root(w http.ResponseWriter, r *http.Request) {
	log.Println("root handler")
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		Render(w, "404page")
		return
	}
	Render(w, "")
}