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

	// * save the file to disk
	if err := SaveFileToDisk(header, file); err != nil {
		log.Println("Error saving file to upload dir:", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// * save the location to the database
	documentID, err := PersistDocumentMetaData(r.Context(), header, file)
	if err != nil {
		log.Println("Error saving document metadata:", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// * go through the db and create or update the records
	if err := PersistTransactions(r.Context(), documentID, header, file); err != nil {
		log.Println("Error saving document metadata:", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	// * respond with the values
	// Log the file name
	log.Printf("Uploaded file: %s", header.Filename)

	// You can add further processing of the file here

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}
