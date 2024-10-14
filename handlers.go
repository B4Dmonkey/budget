package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cbroglie/mustache"
)

func (a *App) addHandlers() {
	handlers := []struct {
		method  string
		path    string
		handler http.HandlerFunc
	}{
		{method: http.MethodGet, path: "/", handler: WithMiddleware(a.root, logger)},
		{method: http.MethodPost, path: "/documents", handler: WithMiddleware(a.documents, logger)},
	}

	for _, h := range handlers {
		a.mux.HandleFunc(fmt.Sprintf("%s %s", h.method, h.path), h.handler)
	}
}

func (a *App) root(w http.ResponseWriter, r *http.Request) {
	homePage := HomePage{ctx: r.Context(), q: a.db_queries}
	binding := homePage.GetBinding()
	template, err := homePage.Template()

	if err != nil {
		log.Println("Error rendering template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)

	if err := template.FRender(w, binding); err != nil {
		log.Println("Error rendering template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func (a *App) documents(w http.ResponseWriter, r *http.Request) {
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

		return
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		http.Error(w, "Error reading file", http.StatusBadRequest)

		return
	}

	// ?This should be programmatically determined
	start_date := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.Local)
	end_date := time.Date(2024, time.September, 30, 23, 59, 59, 999999999, time.Local)

	transactions, err := GetTransactions(a.db_queries, r.Context(), start_date, end_date)
	if err != nil {
		log.Println("Error getting transactions:", err)
		http.Error(w, "Error getting transactions", http.StatusInternalServerError)
	}

	template_file, err := viewsDir.ReadFile("views/components/transactions.mst")
	if err != nil {
		log.Println("Error reading template file:", err)
		http.Error(w, "Error reading template file", http.StatusInternalServerError)
	}

	template, _ := mustache.ParseString(string(template_file))

	var viewableTransactions []TransactionViewModel

	for _, pt := range transactions {
		balance, _ := pt.Transaction.Balance.Value()

		var balanceStr string

		if balance == nil {
			balanceStr = ""
		} else {
			balanceStr = ConvertCurrencyIntToString(balance.(int64))
		}

		viewableTransactions = append(viewableTransactions, TransactionViewModel{
			Date:        pt.Transaction.PostingDate.Format("2006-01-02"),
			Amount:      ConvertCurrencyIntToString(pt.Transaction.Amount),
			Description: pt.Transaction.Description,
			Balance:     balanceStr,
		})
	}

	w.WriteHeader(http.StatusCreated)

	if err := template.FRender(
		w,
		HomePageViewModel{
			HasUnprocessedTransactions: len(viewableTransactions) > 0,
			UnprocessedTransactions:    viewableTransactions,
		}); err != nil {
		log.Println("Error rendering template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
