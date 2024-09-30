package main

import (
	"log"
	"my-budget/database/orm"
	"net/http"
	"time"
)


func root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		Render(w, "404page")
		return
	}
	db := orm.New(conn)
	startDate := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(2024, time.September, 30, 23, 59, 59, 999999999, time.Local)
	pending_transactions, err := db.GetTransactionsInDateRange(r.Context(), orm.GetTransactionsInDateRangeParams{
		PostingDate:   startDate,
		PostingDate_2: endDate,
	})
	if err != nil {
		log.Println("Error getting pending transactions:", err)
	}
	transactions_slice := TransactionSlice(pending_transactions)
	transactions_slice.Render(w)
	// Render(w, "", map[string][]orm.Transaction{"unprocessedTransactions": pending_transactions})
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
