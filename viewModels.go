package main

import (
	"log"
	"time"

	"github.com/cbroglie/mustache"

	"my-budget/app"
	"my-budget/database/orm"
)

type TransactionViewModelSlice []TransactionViewModel

type HomePage struct {
	ctx app.Context
}

type TransactionViewModel struct {
	Date        string
	Amount      string
	Description string
	Balance     string
}

func (h HomePage) Binding() interface{} {
	start_date := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.Local)
	end_date := time.Date(2024, time.September, 30, 23, 59, 59, 999999999, time.Local)
	pending_transactions, err := h.ctx.DB.(*orm.Queries).GetTransactionsInDateRange(
		h.ctx.Req.Context(),
		orm.GetTransactionsInDateRangeParams{PostingDate: start_date,
			PostingDate_2: end_date,
		})

	if err != nil {
		log.Println("Error getting pending transactions:", err)
	}

	var transactionsView TransactionViewModelSlice
	for _, transaction := range pending_transactions {
		balance, _ := transaction.Balance.Value()
		var balanceStr string
		if balance == nil {
			balanceStr = ""
		} else {
			balanceStr = ConvertCurrencyIntToString(balance.(int64))
		}

		transactionsView = append(transactionsView, TransactionViewModel{
			Date:        transaction.PostingDate.Format("2006-01-02"),
			Amount:      ConvertCurrencyIntToString(transaction.Amount),
			Description: transaction.Description,
			Balance:     balanceStr,
		})
	}
	return map[string]TransactionViewModelSlice{"UnprocessedTransactions": transactionsView}
}

func (h HomePage) Template() (*mustache.Template, error) {
	if template_file, err := viewsDir.ReadFile("views/pages/home.mst"); err != nil {
		return nil, err
	} else {
		return mustache.ParseString(string(template_file))
	}
}
