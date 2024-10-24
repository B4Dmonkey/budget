package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
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

func isSpotifyBill(date time.Time) bool {
	lastDayOfMonth := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day()
	if lastDayOfMonth >= 29 && date.Day() == 29 {
		return true
	}

	if lastDayOfMonth == 28 && date.Day() == 28 {
		return true
	}

	return false
}

func newTransactionCSV(file *os.File) *csv.Writer {
	writer := csv.NewWriter(file)
	header := []string{"Details", "Posting Date", "Description", "Amount", "Type", "Balance", "Check or Slip #"}

	if err := writer.Write(header); err != nil {
		log.Fatalf("Failed to write header: %s", err)
	}

	return writer
}

func (t Transaction) transactionToRow() []string {
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

func writeTransactions(writer *csv.Writer, t []Transaction) {
	for i := len(t) - 1; i >= 0; i-- {
		if err := writer.Write(t[i].transactionToRow()); err != nil {
			println("Failed to write row: ", err)
			continue
		}
	}
}

func main() {
	SEED := int64(1)
	r := rand.New(rand.NewSource(SEED))
	println("Seed: ", SEED)
	start_balance := 10 + r.Float64()*(900-10)
	start_balance = math.Round(start_balance*100) / 100
	println(fmt.Sprintf("Starting Balance: %.2f", start_balance))
	start_balance += 3500
	println(fmt.Sprintf("Starting Balance after Payroll: %.2f", start_balance))

	// ! Example of using faker. I may not need it
	fake := faker.NewWithSeed(r)
	fake.Person().Name()
	// ! End of example

	file, err := os.Create("/Users/appstack/Developer/Personal/budget/testdata/Chase Activity Feb 1.CSV")
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	defer file.Close()

	writer := newTransactionCSV(file)
	defer writer.Flush()

	start_date := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	end_date := time.Date(start_date.Year(), start_date.Month()+1, 0, 0, 0, 0, 0, time.UTC)

	balance := start_balance
	var transactions []Transaction
	for date := start_date; !date.After(end_date); date = date.AddDate(0, 0, 1) {
		formattedDateString := date.Format("01/02/2006")

		if isPayDay(date) {
			amount := 3500.00
			balance += amount
			transactions = append(transactions, Transaction{
				Balance:     fmt.Sprintf("%.2f", balance),
				Details:     CREDIT,
				PostingDate: formattedDateString,
				Description: "PAYROLL",
				Amount:      fmt.Sprintf("%.2f", amount),
				Type:        "ACH_CREDIT",
			})
		}

		if isSpotifyBill(date) {
			amount := -11.99
			balance += amount
			transactions = append(transactions, Transaction{
				Balance:     fmt.Sprintf("%.2f", balance),
				Details:     DEBIT,
				PostingDate: formattedDateString,
				Description: "Spotify Bill Proxy",
				Amount:      fmt.Sprintf("%.2f", amount),
				Type:        "DEBIT_CARD",
			})
		}
	}

	writeTransactions(writer, transactions)
}
