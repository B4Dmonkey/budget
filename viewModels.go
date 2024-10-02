package main

import (
	"github.com/cbroglie/mustache"
	"log"
	"my-budget/app"
	"my-budget/database/orm"
	"time"
)

type TransactionSlice []orm.Transaction

type HomePage struct {
	ctx app.Context
}

func (h HomePage) Binding() interface{} {
	startDate := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(2024, time.September, 30, 23, 59, 59, 999999999, time.Local)
	pending_transactions, err := h.ctx.DB.(*orm.Queries).GetTransactionsInDateRange(
		h.ctx.Req.Context(),
		orm.GetTransactionsInDateRangeParams{PostingDate: startDate,
			PostingDate_2: endDate,
		})

	if err != nil {
		log.Println("Error getting pending transactions:", err)
	}

	return map[string]TransactionSlice{"UnprocessedTransactions": pending_transactions}
}
func (h HomePage) Template() (*mustache.Template, error) {
	if template_file, err := viewsDir.ReadFile("views/pages/home.mst"); err != nil {
		return nil, err
	} else {
		return mustache.ParseString(string(template_file))
	}
}
