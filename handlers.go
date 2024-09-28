package main

import (
	"io"
	"log"
	"net/http"
	"os"
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

	// * save the file to disk
	// todo: check if the folder
	out, err := os.Create("uploads/" + header.Filename)
	if err != nil {
		log.Println("Error creating file:", err)
		http.Error(w, "Unable to create the file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Println("Error saving file:", err)
		http.Error(w, "Unable to save the file", http.StatusInternalServerError)
		return
	}
	// * save the location to the database
	// * go through the db and create or update the records
	// * respond with the values
	// Log the file name
	log.Printf("Uploaded file: %s", header.Filename)

	// You can add further processing of the file here

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}