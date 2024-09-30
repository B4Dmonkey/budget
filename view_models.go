package main

import (
	"fmt"
	"log"
	"my-budget/database/orm"
	"net/http"

	"github.com/cbroglie/mustache"
)

type TransactionSlice []orm.Transaction

func (t TransactionSlice) Template() string {
	template, err := viewsDir.ReadFile("views/pages/home.mst")
	if err != nil {
		log.Fatal("Error reading file", err)
	}
	return string(template)
}

func (t TransactionSlice) Render(w http.ResponseWriter) {
	template := t.Template()
	parsedTemplate, _ := mustache.ParseString(template)
	content, err := parsedTemplate.Render(template, map[string]TransactionSlice{"UnprocessedTransactions": t})
	if err != nil {
		log.Fatal("Error rendering template", err)
	}
	fmt.Fprint(w, content)
}
