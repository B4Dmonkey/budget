package main

import (
	"fmt"
	"log"
	"net/http"

	"my-budget/app"
	"my-budget/database/orm"
)

func (a *App) addHandlers() {
	handlers := []struct {
		method  string
		path    string
		handler http.HandlerFunc
	}{
		{method: http.MethodGet, path: "/", handler: a.root},
		// {method: http.MethodPost, path: "/documents", handler: a.documents},
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

	tx, err := conn.Begin()
	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadRequest)
	}
	defer tx.Rollback()
	db := orm.New(conn)
	txdb := db.WithTx(tx)

	if err := ProcessNewTransactions(txdb, r.Context(), header, file); err != nil {
		log.Println("Error processing transactions:", err)
		http.Error(w, "Error reading file", http.StatusBadRequest)
	}

	tx.Commit()
}

// func root(ctx app.Context) error {
// 	log.Println("Root handler")
// 	homePage := HomePage{ctx: ctx}
// 	return ctx.Render(http.StatusOK, homePage)
// }

func documents(ctx app.Context) error {
	log.Println("Documents handler")
	file, header, err := ctx.Req.FormFile("file")
	if err != nil {
		log.Println("Error reading file:", err)
		return err
	}
	defer file.Close()

	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	db := orm.New(conn)
	txdb := db.WithTx(tx)

	if err := ProcessNewTransactions(txdb, ctx.Req.Context(), header, file); err != nil {
		log.Println("Error processing transactions:", err)
		return err
	}

	tx.Commit()
	return nil
}
