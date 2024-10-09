package main

import (
	"fmt"
	"log"
	"net/http"
)

func (a *App) addHandlers() {
	handlers := []struct {
		method  string
		path    string
		handler http.HandlerFunc
	}{
		{method: http.MethodGet, path: "/", handler: a.root},
		{method: http.MethodPost, path: "/documents", handler: a.documents},
	}

	for _, h := range handlers {
		a.mux.HandleFunc(fmt.Sprintf("%s %s", h.method, h.path), h.handler)
	}
}

func (a *App) root(w http.ResponseWriter, r *http.Request) {
	log.Println("Root handler")
	homePage := HomePage{ctx: r.Context(), q: a.db_queries}
	content, err := a.Render(homePage)
	if err != nil {
		log.Println("Error rendering template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(content))
}

func (a *App) documents(w http.ResponseWriter, r *http.Request) {
	log.Println("Documents handler")
	file, header, err := r.FormFile("file")

	if err != nil {
		log.Println("Error reading file:", err)
		http.Error(w, "Error reading file", http.StatusBadRequest)
	}
	defer file.Close()

	tx, err := a.db_conn.Begin()
	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadRequest)
	}
	defer tx.Rollback()
	txdb := a.db_queries.WithTx(tx)

	if err := ProcessNewTransactions(txdb, r.Context(), header, file); err != nil {
		log.Println("Error processing transactions:", err)
		http.Error(w, "Error reading file", http.StatusBadRequest)
	}

	tx.Commit()
}
