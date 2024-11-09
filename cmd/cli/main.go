package main

import (
	"context"
	"log"

	db "my-budget/database"
	// "my-budget/database/orm"
	"my-budget/internal/budget"
	"os"
)

const AddNewTransactionsFromDocument = "AddNewTransactionsFromDocument"
const BulkAddNewTransactionsFromDocument = "BulkAddNewTransactionsFromDocument"
const GetIncomeVsExpense = "GetIncomeVsExpense"
const GetPendingTransactions = "GetPendingTransactions"

// const Request = AddNewTransactionsFromDocument
const Request = BulkAddNewTransactionsFromDocument

func main() {
	ctx := context.Background()

	var ok bool

	var DATABASE_LOC string

	if DATABASE_LOC, ok = os.LookupEnv("DATABASE_LOC"); !ok || DATABASE_LOC == "" {
		log.Fatal(ok, "DATABASE_LOC not set")
	}

	conn := db.Connect(DATABASE_LOC)

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

	case BulkAddNewTransactionsFromDocument:
		base_dir := "/Users/appstack/Developer/Personal/budget/cmd/generateTestData"
		fileNames := []string{
			"Chase Activity Sept 27.CSV",
			"Chase Activity Oct 6.CSV",
			"Chase Activity Oct 30.CSV",
			// "Chase9931_Activity_20240412.CSV",
		}
		transactions := []struct {
			File     *os.File
			FileName string
		}{}

		for _, fileName := range fileNames {
			file, err := os.Open(base_dir + "/" + fileName)
			if err != nil {
				log.Fatal(err)
			}

			transactions = append(transactions, struct {
				File     *os.File
				FileName string
			}{file, fileName})
		}

		for _, transaction := range transactions {
			if err := budget.AddNewTransactionsFromDocument(transaction.FileName, transaction.File); err != nil {
				log.Fatal(err)
			}
		}

		// q := orm.New(conn)

	case GetIncomeVsExpense:
		if err := budget.GetIncomeVsExpense(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Invalid request")
	}
}
