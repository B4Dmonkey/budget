package main

import (
	"context"
	"log"
	"my-budget/internal"
	"os"
)

func main() {
	ctx := context.Background()
	conn, err := internal.ConnectToDatabase(ctx)

	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("/Users/appstack/Developer/Personal/budget/cmd/generateTestData/Chase Activity Oct 6.CSV")
	if err != nil {
		log.Fatal(err)
	}
	
	budget, err := internal.NewBudget(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}

	if err := budget.AddNewTransactionsFromDocument("Chase Activity Oct 6.CSV", file); err != nil {
		log.Fatal(err)
	}
}
