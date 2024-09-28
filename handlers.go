package main

import (
	"log"
	"net/http"
)

func root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		Render(w, "404page")
		return
	}
	Render(w, "")
}

func upload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println("Error reading file:", err)
		http.Error(w, "Unable to read file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Log the file name
	log.Printf("Uploaded file: %s", header.Filename)

	// You can add further processing of the file here

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}