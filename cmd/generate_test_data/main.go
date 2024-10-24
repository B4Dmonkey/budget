package main

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jaswdr/faker/v2"
)

const CREDIT = "CREDIT"
const DEBIT = "DEBIT"

type Transaction struct {
	Details           string
	PostingDate       string
	Description       string
	Amount            string
	Type              string
	Balance           string
	CheckOrSlipNumber string
}

func isPayDay(date time.Time) bool {
	fifteenth := 15
	lastDayOfMonth := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day()
	isWeekend := date.Weekday() == time.Saturday || date.Weekday() == time.Sunday
	isFriday := date.Weekday() == time.Friday

	isFifteenthOrLastDayOfTheMonth := func(d int) bool { return d == fifteenth || d == lastDayOfMonth }

	isWeekendFifteenthOrLastDayOfTheMonth := func() bool {
		if !isFriday {
			return false
		}
		if saturday := date.AddDate(0, 0, 1); isFifteenthOrLastDayOfTheMonth(saturday.Day()) {
			return true
		}
		if sunday := date.AddDate(0, 0, 2); isFifteenthOrLastDayOfTheMonth(sunday.Day()) {
			return true
		}
		return false
	}

	if isFriday && isWeekendFifteenthOrLastDayOfTheMonth() {
		return true
	}

	if !isWeekend && isFifteenthOrLastDayOfTheMonth(date.Day()) {
		return true
	}

	return false
}

func transactionToRow(t Transaction) []string {
	return []string{
		t.Details,
		t.PostingDate,
		t.Description,
		t.Amount,
		t.Type,
		t.Balance,
		t.CheckOrSlipNumber,
	}
}

func newTransactionCSV(file *os.File) *csv.Writer {
	writer := csv.NewWriter(file)
	header := []string{"Details", "Posting Date", "Description", "Amount", "Type", "Balance", "Check or Slip #"}

	if err := writer.Write(header); err != nil {
		log.Fatalf("Failed to write header: %s", err)
	}

	return writer
}

func main() {
	time_seed := time.Now().UnixNano()

	println("time_seed: ", time_seed)
	seed := int64(1)
	fake := faker.NewWithSeed(rand.NewSource(seed))

	println(fake.Person().Name())

	file, err := os.Create("testdata/Chase Activity Feb 1.CSV")
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	defer file.Close()

	newTransactionCSV(file)

	writer := newTransactionCSV(file)
	defer writer.Flush()

	start_date := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	end_date := time.Date(start_date.Year(), start_date.Month()+1, 0, 0, 0, 0, 0, time.UTC)

	var transactions []Transaction
	for date := start_date; !date.After(end_date); date = date.AddDate(0, 0, 1) {
		formattedDateString := date.Format("01/02/2006")

		if isPayDay(date) {
			transactions = append(transactions, Transaction{
				Details:     CREDIT,
				PostingDate: formattedDateString,
				Description: "PAYROLL",
				Amount:      "3500.00",
				Type:        "ACH_CREDIT",
			})
			continue
		}

		transactions = append(transactions, Transaction{
			PostingDate: formattedDateString,
		})
	}

	for i := len(transactions) - 1; i >= 0; i-- {
		if err := writer.Write(transactionToRow(transactions[i])); err != nil {
			println("Failed to write row: ", err)
			continue
		}
	}
}
