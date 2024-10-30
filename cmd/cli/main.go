package main

import (
	"context"
	"log"
	"my-budget/internal/budget"
	"my-budget/internal/db"
	"os"
)

const AddNewTransactionsFromDocument = "AddNewTransactionsFromDocument"
const BulkAddNewTransactionsFromDocument = "BulkAddNewTransactionsFromDocument"
const GetIncomeVsExpense = "GetIncomeVsExpense"

// const Request = AddNewTransactionsFromDocument
const Request = GetIncomeVsExpense

func main() {
	ctx := context.Background()
	conn, err := db.ConnectToDatabase(ctx)

	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("/Users/appstack/Developer/Personal/budget/cmd/generateTestData/Chase Activity Oct 6.CSV")
	if err != nil {
		log.Fatal(err)
	}

	budget, err := budget.NewBudget(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}

	switch Request {
	case AddNewTransactionsFromDocument:
		if err := budget.AddNewTransactionsFromDocument("Chase Activity Oct 6.CSV", file); err != nil {
			log.Fatal(err)
		}

	case GetIncomeVsExpense:
		if err := budget.GetIncomeVsExpense(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Invalid request")
	}
}
